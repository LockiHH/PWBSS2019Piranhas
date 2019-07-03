package main

import (
	"PWBSS2019/gamelogic"
	"encoding/xml"
	"fmt"
	"github.com/pborman/getopt"
	"io"
	"net"
	"os"
)

type Room struct {
	Id string `xml:"roomId"`
}

type Joined struct {
	Id string `xml:"roomId"`
}

type WelcomeMessage struct {
	Color string `xml:"color,attr"`
}

type MementoMessage struct {
	State StateMessage `xml:"state"`
}

type PlayerMessage struct {
	DisplayName string `xml:"displayName,attr"`
	Color       string `xml:"color,attr"`
}

type StateMessage struct {
	RedPlayer  PlayerMessage `xml:"red"`
	BluePlayer PlayerMessage `xml:"blue"`
	Board      BoardMessage  `xml:"board"`
	StartPlayerColor string `xml:"startPlayerColor,attr"`
	CurrentPlayerColor string `xml:"currentPlayerColor,attr"`
	Turn int `xml:"turn,attr"`
	LastMove *MoveMessage `xml:"lastMove"`
}

type MoveMessage struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
	Direction string `xml:"direction,attr"`
}

type FieldMessage struct {
	FieldState string `xml:"state,attr"`
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}

type FieldsMessage struct {
	Fields []FieldMessage `xml:"field"`
}

type BoardMessage struct {
	Fields []FieldsMessage `xml:"fields"`
}

func StringToFieldType(s string) gamelogic.FieldType {
	switch s {
	case "EMPTY":
		return gamelogic.FieldTypeEmpty
	case "RED":
		return gamelogic.FieldTypeRed
	case "BLUE":
		return gamelogic.FieldTypeBlue
	case "OBSTRUCTED":
		return gamelogic.FieldTypeObstructed
	}
	return gamelogic.FieldTypeEmpty
}

func StringToColor(s string) gamelogic.Color {
	switch s {
	case "blue":
		return gamelogic.ColorBlue
	case "red":
		return gamelogic.ColorRed
	}
	return gamelogic.ColorBlue
}

func createGameState(state *StateMessage) *gamelogic.GameState {
	fields := make([][]*gamelogic.Field, 10)
	for i := 0; i < len(fields); i++ {
		fields[i] = make([]*gamelogic.Field, 10)
	}

	for _, f := range state.Board.Fields {

		for _, field := range f.Fields {
			fields[field.Y][field.X] = &gamelogic.Field{X:field.X, Y:field.Y, T:StringToFieldType(field.FieldState)}
		}
	}

	board := gamelogic.NewBoard(fields, 10, 10)

	return gamelogic.NewGameState(board)
}

func Process(r io.Reader, w io.Writer) error {
	d := xml.NewDecoder(r)
	for {
		v, err := d.Token()
		if err != nil {
			return err
		}

		switch t := v.(type) {

		case xml.StartElement:
			switch t.Name.Local {
			case "data":
				var class string = ""
				for _, v := range t.Attr {
					if v.Name.Local == "class" {
						class = v.Value
						break
					}
				}
				switch class {
				case "memento":
					data := new(MementoMessage)
					err := d.DecodeElement(data, &t)
					if err != nil {
						return err
					}
					gamelogic.GetController().UpdateState(createGameState(&data.State))

				case "welcomeMessage":
					data := new(WelcomeMessage)
					err := d.DecodeElement(data, &t)
					if err != nil {
						return err
					}
					gamelogic.GetController().SetPlayer(StringToColor(data.Color))
				case "sc.framework.plugins.protocol.MoveRequest":
					move, err := gamelogic.GetController().NextTurn()
					roomID := gamelogic.GetController().RoomID()
					if err != nil {
						panic(err)
					}
					io.WriteString(w, fmt.Sprintf("<room roomId=\"%s\"><data class=\"move\" x=\"%d\" y=\"%d\" direction=\"%s\" /></room>", roomID, move.X, move.Y, move.Direction.String()))
				default:
					fmt.Printf("got data of class %s\n", class)

				}
			case "joined":
				for _, v := range t.Attr {
					if v.Name.Local == "roomId" {
						gamelogic.GetController().JoinRoom(v.Value)
						break
					}
				}

			default:
				//fmt.Printf("got xml start tag %s\n", t.Name.Local)
			}

		case xml.EndElement:
			//fmt.Printf("got xml end tag %s\n", t.Name.Local)
		}
	}
}

func main() {
	fmt.Println(os.Args)
	testmode := getopt.BoolLong("test", 't', "")
	host := getopt.StringLong("host", 'h', "localhost", "")
	port := getopt.IntLong("port", 'p', 13050, "")
	reservation := getopt.StringLong("reservation", 'r', "", "")
	getopt.Parse()

	if *testmode {
		file, _ :=os.Open("i.xml")
		err := Process(file, os.Stderr)
		if err != nil {
			panic(err)
		}
		return
	}

	con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		panic(fmt.Sprintf("could not connect to server %s:%d", *host, *port))
	}
	fmt.Println("connected to server")
	//	d := xml.NewDecoder(con)
	_, err = io.WriteString(con, "<protocol>")
	if err != nil {
		panic(err)
	}
	if *reservation == "" {
		_, err =io.WriteString(con, "<join gameType=\"swc_2019_piranhas\"/>")
		if err != nil {
			panic(err)
		}
	} else {
		_, err =io.WriteString(con, fmt.Sprintf("<joinPrepared reservationCode=\"%s\"/>", *reservation))
		if err != nil {
			panic(err)
		}
	}
	err = Process(con, con)
	fmt.Printf("Err: Mensch vs AI-Manual%#v\n", err)
}