package twirc

import (
	"fmt"
)

var mvm_tags = map[string][]string{
	"Operation Oil Spill Badge": {
		"Completed: Decoy Doe's Doom",
		"Completed: Decoy Day of Wreckening",
		"Completed: Coal Town Cave-in",
		"Completed: Coal Town Quarry",
		"Completed: Mannworks Mannhunt",
		"Completed: Mannworks Mean Machines",
	},
	"Operation Steel Trap Badge": {
		"Completed: Coal Town Ctrl+Alt+Destruction",
		"Completed: Coal TownCPU Slaughter",
		"Completed: Decoy Disk Deletion",
		"Completed: Decoy Data Demolition",
		"Completed: Mannworks Machine Massacre",
		"Completed: Mannworks Mech Mutilation",
	},
	"Operation Mecha Engine Badge": {
		"Completed: BigRock Bone shaker",
		"Completed: BigRock Broken Parts",
		"Completed: Decoy Disintegration",
	},
	"Operation Gear Grinder Badge": {
		"Completed: Desperation",
		"Completed: Cataclysm",
		"Completed: Mannslaughter",
	},
	"Operation Two Cities Badge": {
		"Completed: Mannhattan Metro Malice",
		"Completed: Manbhattan Empire Escalation",
		"Completed: Rottenburg Bavarian Botbash",
		"Completed: Rottenburg Hamlet Hospitality",
	},
}

type MvMTour struct {
	Name      string
	Tours     uint64
	Missions  []string
	Completed []string
}

func NewMvMTour(name string, tours uint64, missions []string) MvMTour {
	t := MvMTour{
		Name:      name,
		Tours:     tours,
		Missions:  missions,
		Completed: []string{},
	}
	return t
}

func (mvm *MvMTour) AddCompleted(mission string) {
	mvm.Completed = append(mvm.Completed, mission)
}

func (mvm *MvMTour) InfoStr() string {
	return fmt.Sprintf("%s(%d): %d/%d", mvm.ShortName(), mvm.Tours, len(mvm.Completed), len(mvm.Missions))
}

func (mvm *MvMTour) ShortName() string {
	name := "unknown"
	if mvm.Name == "Operation Oil Spill Badge" {
		name = "OilSpill"
	} else if mvm.Name == "Operation Steel Trap Badge" {
		name = "SteelTrap"
	} else if mvm.Name == "Operation Mecha Engine Badge" {
		name = "Mecha"
	} else if mvm.Name == "Operation Gear Grinder Badge" {
		name = "GearGrinder"
	} else if mvm.Name == "Operation Two Cities Badge" {
		name = "TwoCities"
	}
	return name
}
