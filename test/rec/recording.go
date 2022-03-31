package rec

var recording = make([]string, 0)

func Reset() {
	recording = make([]string, 0)
}

func Add(entry string) {
	recording = append(recording, entry)
}

func Get() []string {
	return recording
}
