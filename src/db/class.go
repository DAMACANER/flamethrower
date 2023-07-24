package db

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
)

type ClassListColumns struct {
	ID                   int            `db:"id"`
	Name                 string         `db:"name"`
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
