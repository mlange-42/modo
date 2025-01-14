package document

type Processor struct {
	Formatter Formatter
}

func NewProcessor(f Formatter) Processor {
	return Processor{
		Formatter: f,
	}
}
