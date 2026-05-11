package scheduling

import "sort"

type PreviewContext struct {
	Variables         []Variable
	PresetConflicts   []PreviewConflict
	Domains           map[string][]DomainValue
	NoDomainConflicts map[string]PreviewConflict
	CandidateOptions  map[string][]PreviewCandidateOption
}

func BuildPreviewContext(input SolverInput) PreviewContext {
	problem := prepareSchedulingProblem(input)
	manualDomains := buildManualOverrideDomains(problem.variables, input)
	mergedDomains := mergeDomainSets(problem.domains, manualDomains)

	return PreviewContext{
		Variables:         problem.variables,
		PresetConflicts:   problem.presetConflicts,
		Domains:           mergedDomains,
		NoDomainConflicts: problem.noDomainConflicts,
		CandidateOptions:  buildCandidateOptions(problem.variables, mergedDomains),
	}
}

func buildCandidateOptions(variables []Variable, domains map[string][]DomainValue) map[string][]PreviewCandidateOption {
	options := make(map[string][]PreviewCandidateOption, len(variables))
	for _, variable := range variables {
		values := domains[variable.ID]
		if len(values) == 0 {
			options[variable.ID] = []PreviewCandidateOption{}
			continue
		}

		items := make([]PreviewCandidateOption, 0, len(values))
		for _, value := range values {
			items = append(items, PreviewCandidateOption{
				Key:          previewCandidateKey(value),
				RoomID:       value.RoomID,
				RoomName:     value.RoomName,
				RoomCapacity: value.RoomCapacity,
				ShiftID:      value.TimeSlot.ShiftID,
				ShiftCode:    value.TimeSlot.ShiftCode,
				ShiftName:    value.TimeSlot.ShiftName,
				ShiftType:    value.TimeSlot.ShiftType,
				StartTime:    value.TimeSlot.Start,
				EndTime:      value.TimeSlot.End,
			})
		}

		sort.Slice(items, func(i, j int) bool {
			if items[i].StartTime.Equal(items[j].StartTime) {
				return items[i].RoomName < items[j].RoomName
			}
			return items[i].StartTime.Before(items[j].StartTime)
		})
		options[variable.ID] = items
	}

	return options
}

func previewCandidateKey(value DomainValue) string {
	return slotKeyFromDomain(value) + "|" + value.RoomID
}

func buildManualOverrideDomains(variables []Variable, input SolverInput) map[string][]DomainValue {
	domains := make(map[string][]DomainValue, len(variables))
	for _, variable := range variables {
		slots := generateTimeSlots(input.DateFrom, input.DateTo, variable.DurationMinutes, input.Shifts)
		if len(slots) == 0 {
			domains[variable.ID] = []DomainValue{}
			continue
		}

		values := make([]DomainValue, 0)
		for _, room := range input.Rooms {
			if room.Capacity < variable.ExpectedCapcity {
				continue
			}

			for _, slot := range slots {
				values = append(values, DomainValue{
					RoomID:       room.ID,
					RoomName:     room.Name,
					RoomCapacity: room.Capacity,
					TimeSlot:     slot,
				})
			}
		}

		sort.Slice(values, func(i, j int) bool {
			if values[i].TimeSlot.Start.Equal(values[j].TimeSlot.Start) {
				if values[i].RoomID == variable.PreferredRoomID && values[j].RoomID != variable.PreferredRoomID {
					return true
				}
				if values[j].RoomID == variable.PreferredRoomID && values[i].RoomID != variable.PreferredRoomID {
					return false
				}
				return values[i].RoomName < values[j].RoomName
			}
			return values[i].TimeSlot.Start.Before(values[j].TimeSlot.Start)
		})

		domains[variable.ID] = values
	}

	return domains
}

func mergeDomainSets(primary map[string][]DomainValue, secondary map[string][]DomainValue) map[string][]DomainValue {
	merged := make(map[string][]DomainValue, len(primary)+len(secondary))
	for variableID, values := range primary {
		merged[variableID] = append([]DomainValue(nil), values...)
	}

	for variableID, values := range secondary {
		existing := merged[variableID]
		if len(existing) == 0 {
			merged[variableID] = append([]DomainValue(nil), values...)
			continue
		}

		seen := make(map[string]struct{}, len(existing))
		for _, value := range existing {
			seen[previewCandidateKey(value)] = struct{}{}
		}

		for _, value := range values {
			key := previewCandidateKey(value)
			if _, ok := seen[key]; ok {
				continue
			}
			existing = append(existing, value)
			seen[key] = struct{}{}
		}

		sort.Slice(existing, func(i, j int) bool {
			if existing[i].TimeSlot.Start.Equal(existing[j].TimeSlot.Start) {
				return existing[i].RoomName < existing[j].RoomName
			}
			return existing[i].TimeSlot.Start.Before(existing[j].TimeSlot.Start)
		})
		merged[variableID] = existing
	}

	return merged
}
