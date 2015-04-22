package main

type Loader struct {
	numOfClients  int
	numOfMessages int
}

func NewLoader(numOfClients int, numOfMessages int) *Loader {
	return &Loader{numOfClients: numOfClients, numOfMessages: numOfMessages}
}

func (l *Loader) Start() {
	l.setup()
	l.loadClients()
}

func (l *Loader) Stop() {
	l.tearDown()
}

func (l *Loader) setup() {

}

func (l *Loader) tearDown() {

}

func (l *Loader) loadClients() {

}
