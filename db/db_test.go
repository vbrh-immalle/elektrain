package db

import (
	"testing"
)

func createTestDb() *HardCodedDb {
	db := &HardCodedDb{}
	db.Load([]Opgave{
		{
			Vraag:    "Testvraag01",
			Antwoord: "Testantwoord01",
		},
		{
			Vraag:    "Testvraag02",
			Antwoord: "Testantwoord02",
		},
	})
	return db
}

func TestCount(t *testing.T) {
	db := createTestDb()

	count := db.Count()
	want := 2
	if count != want {
		t.Fatalf("Returned [%d] while we wanted [%d].\n", count, want)
	}
}

func TestNogJuistTeBeantwoordenOpgaves(t *testing.T) {
	db := createTestDb()

	// we assume order is kept
	antwoord := db.NogJuistTeBeantwoordenOpgaves()[0].Antwoord
	want := "Testantwoord01"
	if antwoord != want {
		t.Fatalf("Returned [%s] while we wanted [%s].\n", antwoord, want)
	}
}

func TestRandomNogJuistTeBeantwoordenOpgaveBoundaries(t *testing.T) {
	// A fresh db should always returs a valid Opgave
	for i := 0; i < 1000; i++ {
		db := createTestDb()
		o, _ := db.RandomNogJuistTeBeantwoordenOpgave()
		if o.Antwoord == "" {
			t.Fatalf("Got an empty Opgave as first Random Opgave.\n")
		}
	}
}

func TestUpdateOpgave(t *testing.T) {
	db := createTestDb()

	o := db.opgaves[0]
	o.AantalFouteAntwoorden += 1
	db.UpdateOpgave(o)

	if db.opgaves[0].AantalFouteAntwoorden != 1 {
		t.Fatalf("Should Update. Expected 1, got [%d]\n", db.opgaves[0].AantalFouteAntwoorden)
	}

}

func TestNogJuistTeBeantwoordenOpgaveShrinking(t *testing.T) {
	db := createTestDb()

	opgs1 := db.NogJuistTeBeantwoordenOpgaves()
	if len(opgs1) != 2 {
		t.Fatalf("Expected 2, got [%d]\n", len(opgs1))
	}
	db.opgaves[0].AantalJuisteAntwoorden += 1

	opgs2 := db.NogJuistTeBeantwoordenOpgaves()
	if len(opgs2) != 1 {
		t.Fatalf("Expected 1, got [%d]\n", len(opgs2))
	}
	db.opgaves[1].AantalJuisteAntwoorden += 1

	opgs3 := db.NogJuistTeBeantwoordenOpgaves()
	if len(opgs3) != 0 {
		t.Fatalf("Expected 0, got [%d]\n", len(opgs3))
	}
}

func TestRandomNogJuistTeBeantwoordenOpgaveFinished(t *testing.T) {
	db := createTestDb()

	// After 2 questions answered, we should get finished-signal
	o1, b1 := db.RandomNogJuistTeBeantwoordenOpgave()
	if !b1 {
		t.Fatalf("First question fail")
	}
	o1.AantalJuisteAntwoorden += 1
	db.UpdateOpgave(o1)

	o2, b2 := db.RandomNogJuistTeBeantwoordenOpgave()
	if !b2 {
		t.Fatalf("Second question fail")
	}
	o2.AantalJuisteAntwoorden += 1
	db.UpdateOpgave(o2)

	o3, b3 := db.RandomNogJuistTeBeantwoordenOpgave()
	if b3 {
		t.Fatalf("hasNext should return false for third question")
	}
	if o3.Vraag != "" {
		t.Fatalf("Third Opgave should be empty but is [%s]\n", o3.Vraag)
	}

}
