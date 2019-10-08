package generic

//     +--------+
// a -->        |
//     |   &&   >---- c=a+b, d=a-b
// b -->        |
//     +--------+

type Box interface{}

type Channel struct {
	ch          chan Box
	typ         string
	name        string
	isReadable  bool
	isWriteable bool
}

type Cell struct {
	typ string
	chs []Channel
	fn  func()
}

func registerChannel(c Channel) Channel {
	return c
}

func registerCell(typ string, register map[string](func() Cell)) {

	// Use closure instead of .. register in hashtable
	var foo = func() Cell {

		// register channels
		channelA := registerChannel(Channel{typ: "int", name: "a", isWriteable: true, isReadable: true})
		channelB := registerChannel(Channel{typ: "int", name: "a", isWriteable: true, isReadable: true})
		channelC := registerChannel(Channel{typ: "int", name: "a", isWriteable: true, isReadable: false})
		channelD := registerChannel(Channel{typ: "int", name: "a", isWriteable: true, isReadable: false})

		// register calc fn
		calcFn := func() {

			valA := <-channelA.ch
			valB := <-channelB.ch

			valC := valA.(int) + valB.(int)
			valD := valA.(int) - valB.(int)

			if channelC.isWriteable {
				channelC.ch <- valC
			}

			if channelD.isWriteable {
				channelD.ch <- valD
			}
		}
		c := Cell{typ: typ, fn: calcFn}
		return c
	}

	register[typ] = foo
}
