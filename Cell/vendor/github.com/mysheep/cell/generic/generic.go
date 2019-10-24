package generic

// Inspired by Alias Wavefront Maya way of creating connectable components

//     +--------+
// a -->        |
//     |   &&   >---- c=a+b, d=a-b
// b -->        |
//     +--------+

type Any interface{}

type Channel struct {
	ch          chan Any
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
		channelB := registerChannel(Channel{typ: "int", name: "b", isWriteable: true, isReadable: true})
		channelC := registerChannel(Channel{typ: "int", name: "c", isWriteable: true, isReadable: false})
		channelD := registerChannel(Channel{typ: "int", name: "d", isWriteable: true, isReadable: false})

		// register calc fn
		// calcFn is a function that calcs output by given input
		calcFn := func() {

			valA := <-channelA.ch
			valB := <-channelB.ch

			// TODO: Convert by given type string
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
