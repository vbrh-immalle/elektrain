package main

import (
	"fmt"
	"math/rand"
)

type ElekTrainDb interface {
	Init()
	Load([]Opgave) // mainly for testing
	Print()        // mainly for testing
	Count() int
	CountBeantwoord() int
	CountJuistBeantwoord() int
	CountTotaalAantalAntwoorden() int
	CountOpgavesTeGaan() int
	NogJuistTeBeantwoordenOpgaves() []Opgave
	RandomNogJuistTeBeantwoordenOpgave() (Opgave, bool)
	CorrectAntwoord(opgave Opgave) string
	UpdateOpgave(opgave Opgave)
}

type Opgave struct {
	Id                     int
	Vraag                  string
	Antwoord               string
	AantalFouteAntwoorden  int
	AantalJuisteAntwoorden int
}

type HardCodedDb struct {
	opgaves []Opgave
}

func (db *HardCodedDb) Init() {
	db.opgaves = []Opgave{
		{
			Vraag:    "Wat is de eenheid van spanning (voluit)?",
			Antwoord: "Volt",
		},
		{
			Vraag:    "Wat is de eenheid van spanning (afkorting)?",
			Antwoord: "V",
		},
		{
			Vraag:    "Wat is de eenheid van elektrische stroom (afkorting)?",
			Antwoord: "A",
		},
		{
			Vraag:    "Wat is de eenheid van elektrische stroom (voluit)?",
			Antwoord: "Ampère",
		},
		{
			Vraag:    "Van welke grootheid is Ω de afkorting?",
			Antwoord: "weerstand",
		},
		{
			Vraag:    "Welke eenheid stelt eigenlijk Joule per seconde (J/s) voor? (afkorting)",
			Antwoord: "W",
		},
		{
			Vraag:    "Wat is de S.I.-eenheid voor hoeveelheid elektrische lading (afkorting)?",
			Antwoord: "C",
		},
		{
			Vraag:    "Wat is de S.I.-eenheid voor hoeveelheid elektrische lading (voluit)?",
			Antwoord: "Coulumb",
		},
		{
			Vraag:    "Welke grootheid krijg je als je een spanning en een stroom vermenigvuldigt?",
			Antwoord: "vermogen",
		},
		{
			Vraag:    "Welke grootheid krijg je als je een spanning en een stroom vermenigvuldigt?",
			Antwoord: "vermogen",
		},
		{
			Vraag:    "Welke (samengestelde) eenheid gebruiken energieleveranciers om te bepalen hoeveel energie je verbruikt hebt?",
			Antwoord: "kWh",
		},
		{
			Vraag:    "Hoe schrijf je de (samengestelde) eenheid kiloWatt-uur?",
			Antwoord: "kWh",
		},
		{
			Vraag:    "Welke eenheid gebruiken we voor het aantal ladingen dat per seconde ergens doorstroomt?",
			Antwoord: "A",
		},
	}
	for i := range db.opgaves {
		db.opgaves[i].Id = i
		db.opgaves[i].AantalJuisteAntwoorden = 0
		db.opgaves[i].AantalFouteAntwoorden = 0
	}
}

func (db *HardCodedDb) Print() {
	for i, o := range db.opgaves {
		fmt.Printf("%d: %v+\n", i, o)
	}
}

func (db *HardCodedDb) Load(opgaves []Opgave) {
	db.opgaves = append(db.opgaves, opgaves...)
	for i := range db.opgaves {
		db.opgaves[i].Id = i
		db.opgaves[i].AantalJuisteAntwoorden = 0
		db.opgaves[i].AantalFouteAntwoorden = 0
	}
}

func (db *HardCodedDb) Count() int {
	return len(db.opgaves)
}

func (db *HardCodedDb) CountBeantwoord() int {
	aantal := 0
	for _, opg := range db.opgaves {
		if opg.AantalFouteAntwoorden > 0 || opg.AantalJuisteAntwoorden > 0 {
			aantal += 1
		}
	}
	return aantal
}

func (db *HardCodedDb) CountJuistBeantwoord() int {
	aantal := 0
	for _, opg := range db.opgaves {
		if opg.AantalJuisteAntwoorden > 0 {
			aantal += 1
		}
	}
	return aantal
}

func (db *HardCodedDb) CountTotaalAantalAntwoorden() int {
	aantal := 0
	for _, opg := range db.opgaves {
		aantal += opg.AantalFouteAntwoorden
		aantal += opg.AantalJuisteAntwoorden
	}
	return aantal
}

func (db *HardCodedDb) CountOpgavesTeGaan() int {
	return db.Count() - db.CountJuistBeantwoord()
}

func (db *HardCodedDb) NogJuistTeBeantwoordenOpgaves() []Opgave {
	var opgaves []Opgave
	for _, opg := range db.opgaves {
		if opg.AantalJuisteAntwoorden == 0 {
			opgaves = append(opgaves, opg)
		}
	}
	return opgaves
}

func (db *HardCodedDb) RandomNogJuistTeBeantwoordenOpgave() (Opgave, bool) {
	opgs := db.NogJuistTeBeantwoordenOpgaves()

	if len(opgs) == 0 {
		return Opgave{}, false
	}

	n := rand.Intn(len(opgs))
	return opgs[n], true
}

func (db *HardCodedDb) CorrectAntwoord(opgave Opgave) string {
	return db.opgaves[opgave.Id].Antwoord
}

func (db *HardCodedDb) UpdateOpgave(opgave Opgave) {
	db.opgaves[opgave.Id] = opgave
}
