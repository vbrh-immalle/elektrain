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
			Vraag:    "Welke grootheid drukken we uit in Ω?",
			Antwoord: "weerstand",
		},
		{
			Vraag:    "Welke eenheid is equivalent met Joule per seconde (J/s) (afkorting)?",
			Antwoord: "W",
		},
		{
			Vraag:    "Wat is de S.I.-eenheid voor hoeveelheid elektrische lading (afkorting)?",
			Antwoord: "C",
		},
		{
			Vraag:    "Wat is de S.I.-eenheid voor hoeveelheid elektrische lading (voluit)?",
			Antwoord: "Coulomb",
		},
		{
			Vraag:    "Welke grootheid krijg je als je een spanning en een stroom vermenigvuldigt?",
			Antwoord: "vermogen",
		},
		{
			Vraag:    "Hoeveel Coulomb (of Ampère-seconde) aan lading wordt verplaatst als een lege batterij van 1000 mAh helemaal wordt opgeladen (getal)?",
			Antwoord: "3600",
		},
		{
			Vraag:    "Welke (samengestelde) eenheid gebruiken energieleveranciers om te bepalen hoeveel energie je verbruikt hebt (afkorting)?",
			Antwoord: "kWh",
		},
		{
			Vraag:    "Hoe schrijf je de (samengestelde) eenheid kiloWatt-uur (afkorting)?",
			Antwoord: "kWh",
		},
		{
			Vraag:    "Welke eenheid gebruiken we voor het aantal ladingen dat per seconde ergens doorstroomt (afkorting)?",
			Antwoord: "A",
		},
		{
			Vraag:    "Hoeveel spanning staat er over een normale USB-poort (getal) (zonder USB-PD (Power Delivery))?",
			Antwoord: "5",
		},
		{
			Vraag:    "Hoeveel seconden zitten er in 1 uur (getal)?",
			Antwoord: "3600",
		},
		{
			Vraag:    "Hoe groot is de netspanning in Europa (getal + eenheid)?",
			Antwoord: "230 V",
		},
		{
			Vraag:    "Hoeveel uur kan je een toestel van 500 W op laten staan als je maar 1 kWh aan energie mag verbuiken (getal)?",
			Antwoord: "2",
		},
		{
			Vraag:    "Welke grootheid krijg je als je een hoeveelheid lading ([Q] of [mAh] of [As]) vermenigvuldigt met een spanning ([V])? TIP: de eenheid is o.a. [kWh]",
			Antwoord: "energie",
		},
		{
			Vraag:    "Wat is de frequentie van de sinusgolf v.d. wisselspanning en -stroomop het Europese elektriciteitsnet(waarde + eenheid)?",
			Antwoord: "50 Hz",
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
