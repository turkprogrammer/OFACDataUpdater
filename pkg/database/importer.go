package database

import "OFACDataUpdater/pkg/model"

// Фильтрация записей по sdnType=Individual
func filterAndImportIndividuals(sdnList model.SDNList) error {
	individuals, err := filterIndividuals(sdnList)
	if err != nil {
		return err
	}

	return UpdateData(individuals)
}
