package pri

type Priority struct {
	Facility Facility
	Severity Severity
}

func FromCombined(p byte) (Priority, error) {
	if e := Facility(p).Masked().Valid(); e != nil {
		return Priority{}, e
	}
	return Priority{
		Facility(p).Masked(),
		Severity(p).Masked(),
	}, nil
}

func (p Priority) Combine() (byte, error) {
	if e := p.Facility.Valid(); e != nil {
		return 0x00, e
	}
	return byte(p.Facility.Masked()) | byte(p.Severity.Masked()), nil
}
