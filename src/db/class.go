package db

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v6"
)

type ClassListColumns struct {
	ID                   sql.NullInt16  `db:"id"`
	Name                 sql.NullString `db:"name"`
	Type                 sql.NullString `db:"type"`
	Alignment            sql.NullString `db:"alignment"`
	HitDie               sql.NullString `db:"hit_die"`
	ClassSkills          sql.NullString `db:"class_skills"`
	SkillPoints          sql.NullString `db:"skill_points"`
	SkillPointsAbility   sql.NullString `db:"skill_points_ability"`
	SpellStat            sql.NullString `db:"spell_stat"`
	Proficiencies        sql.NullString `db:"proficiencies"`
	SpellType            sql.NullString `db:"spell_type"`
	EpicFeatBaseLevel    sql.NullString `db:"epic_feat_base_level"`
	EpicFeatInterval     sql.NullString `db:"epic_feat_interval"`
	EpicFeatList         sql.NullString `db:"epic_feat_list"`
	EpicFullText         sql.NullString `db:"epic_full_text"`
	ReqRace              sql.NullString `db:"req_race"`
	ReqWeaponProficiency sql.NullString `db:"req_weapon_proficiency"`
	ReqBaseAttackBonus   sql.NullString `db:"req_base_attack_bonus"`
	ReqSkill             sql.NullString `db:"req_skill"`
	ReqFeat              sql.NullString `db:"req_feat"`
	ReqSpells            sql.NullString `db:"req_spells"`
	ReqLanguages         sql.NullString `db:"req_languages"`
	ReqPsionics          sql.NullString `db:"req_psionics"`
	ReqEpicFeat          sql.NullString `db:"req_epic_feat"`
	ReqSpecial           sql.NullString `db:"req_special"`
	SpellList1           sql.NullString `db:"spell_list_1"`
	SpellList2           sql.NullString `db:"spell_list_2"`
	SpellList3           sql.NullString `db:"spell_list_3"`
	SpellList4           sql.NullString `db:"spell_list_4"`
	SpellList5           sql.NullString `db:"spell_list_5"`
	FullText             sql.NullString `db:"full_text"`
	Reference            sql.NullString `db:"reference"`
}

var ClassListTableName = "class"

type ClassListRepo struct {
	*BaseRepo
	Class ClassListColumns
}

var TotalClassCount int

func (b *BaseRepo) Count() *BaseRepo {
	b.QueryBuilder = squirrel.Select("COUNT(*)").From(ClassListTableName)
	return b
}

func (c *ClassListRepo) DropTableIfExists() {
	_, err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", ClassListTableName))
	if err != nil {
		log.Fatal(err)
	}
}
func (c *ClassListRepo) CreateTable() {
	// create table from ClassListColumns type
	_, err := DB.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		type TEXT,
		alignment TEXT,
		hit_die TEXT,
		class_skills TEXT,
		skill_points TEXT,
		skill_points_ability TEXT,
		spell_stat TEXT,
		proficiencies TEXT,
		spell_type TEXT,
		epic_feat_base_level TEXT,
		epic_feat_interval TEXT,
		epic_feat_list TEXT,
		epic_full_text TEXT,
		req_race TEXT,
		req_weapon_proficiency TEXT,
		req_base_attack_bonus TEXT,
		req_skill TEXT,
		req_feat TEXT,
		req_spells TEXT,
		req_languages TEXT,
		req_psionics TEXT,
		req_epic_feat TEXT,
		req_special TEXT,
		spell_list_1 TEXT,
		spell_list_2 TEXT,
		spell_list_3 TEXT,
		spell_list_4 TEXT,
		spell_list_5 TEXT,
		full_text TEXT,
		reference TEXT
	)`, ClassListTableName))
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ClassListRepo) PopulateTable(populateSize uint16) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// generate random data with gofakeit and insert em
	for i := 0; i < int(populateSize); i++ {
		c.Class = ClassListColumns{}
		c.Class.Name = sql.NullString{String: gofakeit.Person().FirstName, Valid: true}
		c.Class.Type = sql.NullString{String: gofakeit.Car().Type, Valid: true}
		c.Class.Alignment = sql.NullString{String: gofakeit.BeerStyle(), Valid: true}
		c.Class.HitDie = sql.NullString{String: "1d6", Valid: true}
		c.Class.ClassSkills = sql.NullString{String: gofakeit.HipsterSentence(15), Valid: true}
		c.Class.SkillPoints = sql.NullString{String: fmt.Sprintf("%d", gofakeit.IntRange(1, 15)), Valid: true}
		c.Class.SkillPointsAbility = sql.NullString{String: gofakeit.Bird(), Valid: true}
		c.Class.SpellStat = sql.NullString{String: gofakeit.HackerAbbreviation(), Valid: true}
		c.Class.Proficiencies = sql.NullString{String: gofakeit.Name(), Valid: true}
		c.Class.SpellType = sql.NullString{String: gofakeit.BeerName(), Valid: true}
		c.Class.EpicFeatBaseLevel = sql.NullString{String: gofakeit.Name(), Valid: true}
		c.Class.EpicFeatInterval = sql.NullString{String: gofakeit.Name(), Valid: true}
		c.Class.EpicFeatList = sql.NullString{String: gofakeit.Name(), Valid: true}
		c.Class.EpicFullText = sql.NullString{String: gofakeit.Name(), Valid: true}
		c.Class.ReqRace = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqWeaponProficiency = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqBaseAttackBonus = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqSkill = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqFeat = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqSpells = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqLanguages = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqPsionics = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqEpicFeat = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.ReqSpecial = sql.NullString{String: gofakeit.HipsterSentence(20), Valid: true}
		c.Class.SpellList1 = sql.NullString{String: gofakeit.HipsterWord(), Valid: true}
		c.Class.SpellList2 = sql.NullString{String: gofakeit.HipsterWord(), Valid: true}
		c.Class.SpellList3 = sql.NullString{String: gofakeit.HipsterWord(), Valid: true}
		c.Class.SpellList4 = sql.NullString{String: gofakeit.HipsterWord(), Valid: true}
		c.Class.SpellList5 = sql.NullString{String: gofakeit.HipsterWord(), Valid: true}
		c.Class.FullText = sql.NullString{String: gofakeit.HipsterParagraph(4, 5, 200, "\n"), Valid: true}
		c.Class.Reference = sql.NullString{String: gofakeit.NewCrypto().BuzzWord(), Valid: true}
		fields := c.ExtractVars()
		columns := make([]string, 0, len(fields))
		rows := make([]interface{}, 0, len(fields))
		for k, v := range fields {
			columns = append(columns, k)
			rows = append(rows, v)
		}
		sql, args, err := squirrel.Insert(ClassListTableName).Columns(columns...).Values(rows...).ToSql()
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(err)
		}
		_, err = tx.Exec(sql, args...)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatal(err)
			}

			log.Fatal(err)
		}

	}
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

}

// ExtractVars returns any non-null DB fields with the associated
// Class table inside the ClassRepo.
//
// Formatting is like follows:
//
//	db_fields["reference"] = "example"
//
// Map keys are extracted from structs db tags.
func (c *ClassListRepo) ExtractVars() map[string]interface{} {
	v := reflect.ValueOf(c.Class)
	t := v.Type()
	var fields = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("db")
		switch field.Interface().(type) {
		case sql.NullString:
			ns := field.Interface().(sql.NullString)
			if ns.Valid {
				fields[tag] = ns.String
			}
		case string:
			if field.String() != "" {
				fields[tag] = field.String()
			}
		case int:
			if field.Int() != 0 {
				fields[tag] = fmt.Sprintf("%d", field.Int())
			}
		}
	}
	return fields
}

// Find selects everything from the constan ClassListTableName.
//
// Usage:
//
//	repo := &db.ClassListRepo{BaseRepo: &db.BaseRepo{}}
//	rows, err := repo.Find().Paginate(1, 6).Query()
func (c *ClassListRepo) Find() *BaseRepo {
	c.QueryBuilder = squirrel.Select("*").From(ClassListTableName)
	return c.BaseRepo
}
