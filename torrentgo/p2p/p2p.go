package p2p

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/connection"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/message"
	"github.com/zemberdotnet/gotorrent/piece"
	"github.com/zemberdotnet/gotorrent/state"
	"github.com/zemberdotnet/gotorrent/work"
)

type P2P struct {
	state        *state.TorrentState
	connPool     *connection.ConnectionPool
	workCreator  *work.WorkCreator
	outChan      chan *piece.TorrPiece
	maxBacklog   int
	maxBlocksize int
}

func (p *P2P) Start(ctx context.Context) {

	// Must return a connection? Should we enfoce that
	// or have strategy determine when to throw away
	var conn *connection.PeerConn
	var err error

	// Do we need to defer connection close
	for {
		conn, err = p.connPool.NextConnection()
		if err != nil {
			if err.Error() == "no remaining peers" {
				return
			}
			log.Println(err)
			continue
		}

		//defer p.connPool.ReturnConnection(conn)

		// we initialize the connection
		// complete the handshake
		// and set the bitfield
		err = conn.Initialize()
		if err != nil {
			continue
		}
		log.Println("FINISHED INIT")
		break
	}

	// update the state
	// this is more precarious than I would like
	p.state.IncrementCounts(conn.Bitfield)

	//
	// c.Send(
	//

	// wait for work to come to us
	for {
		select {
		case <-ctx.Done():
			return
		default:
			work, err := p.workCreator.WorkFromBitfield(conn.Bitfield, false)
			if err != nil {
				log.Println("Unable to get work")
				return
			}
			err = p.DownloadPiece(conn, work[0])
			if err != nil {
				log.Printf("TODO: err %v", err)
				return
			}

			if work[0].Verify() != nil {

			}

			// this is an interestin problem to understand
			// calling this first garuntees we don't overflow
			// the buffer
		}
	}

	// need to initiate a new connection
}

// DownloadPiece attempts to downlad a torreent piece
func (p *P2P) DownloadPiece(conn *connection.PeerConn, piece *piece.TorrPiece) (err error) {

	fmt.Println("Downloading Piece")
	// should move this onto an object so we
	// can track the state
	backlog := 0

	// Deadline removes bad connections
	conn.Conn.SetDeadline(time.Now().Add(time.Second * 30))
	defer conn.Conn.SetDeadline(time.Time{})

	for piece.Downloaded < piece.Length() {
		if !conn.Choked {
			// fill the backlog
			for backlog < p.maxBacklog && piece.Requested < piece.Length() {
				blockSize := p.maxBlocksize

				// the last block is often shorter than the full maxBlocksize
				if piece.Length()-piece.Requested < blockSize {
					blockSize = piece.Length() - piece.Requested
				}
				// TODO
				// do checking on blocksize and set on piece
				err = message.SendMessage(conn.Conn, &message.Message{
					MessageID: message.MsgRequest,
					Index:     piece.Index(),
					Begin:     piece.Begin() + piece.Requested,
					Length:    blockSize,
					Payload:   []byte{},
				})

				if err != nil {
					return err
				}

				piece.Requested += p.maxBlocksize
				backlog++
			}

			err = p.handleResponse(conn, piece)
			if err != nil {
				return err
			}

		}

		fmt.Println("Finished Downloading Piece")
	}
	return nil
}

func (p *P2P) handleResponse(conn *connection.PeerConn, piece *piece.TorrPiece) error {

	msg, err := message.ReadMessage(conn.Conn)
	if err != nil {
		return err
	}

	switch msg.MessageID {
	case message.MsgKeepAlive:
		return nil
	case message.MsgChoke:
		conn.Choked = true
	case message.MsgUnchoke:
		conn.Choked = false
	case message.MsgInterested:
		conn.Interested = true
		// TODO
	case message.MsgNotInterested:
		conn.Interested = false
		// TODO
	case message.MsgHave:
		// we are relying on not recieving bad input
		conn.Bitfield.SetPiece(msg.Index)
	case message.MsgBitfield:
		log.Println("TODO CHECK THIS AGAINST SPEC")
		conn.Bitfield = bitfield.NewBitfieldFromBytes(msg.Payload)
	case message.MsgRequest:
		// TODO
	case message.MsgPiece:
		log.Println("Getting a piece!")
		n, err := piece.WriteAt(msg.Payload, msg.Begin)
		if err != nil {
			return err
		}

		if n != len(msg.Payload) {
			log.Printf("TODO: n %v len(msg.Payload) %v", n, len(msg.Payload))
		}
		piece.Downloaded += n
		fmt.Printf("Downloaded %v", n)

		return nil
	case message.MsgCancel:
		log.Println("Cancel Message")
		// TODO
	case message.MsgPort:
		log.Println("Message Port")
		// TODO

	}
	return nil
}

func (p *P2P) verifyPiece(piece *piece.TorrPiece) error {
	return nil
}

func NewP2PFactory(state *state.TorrentState, cp *connection.ConnectionPool, wc *work.WorkCreator, oc chan *piece.TorrPiece) func() interfaces.Strategy {
	return func() interfaces.Strategy {
		return &P2P{
			state:       state,
			connPool:    cp,
			workCreator: wc,
			outChan:     oc,
			// TODO make this configurable
			maxBacklog:   5,
			maxBlocksize: 16384,
		}
	}
}
