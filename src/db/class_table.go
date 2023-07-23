package db

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
)

type ClassTableColumns struct {
	ID                  int            `db:"id"`
	Name                string         `db:"name"`
	Level               sql.NullString `db:"level"`
	BaseAttackBonus     sql.NullString `db:"base_attack_bonus"`
	FortSave            sql.NullString `db:"fort_save"`
	RefSave             sql.NullString `db:"ref_save"`
	WillSave            sql.NullString `db:"will_save"`
	CasterLevel         sql.NullString `db:"caster_level"`
	PointsPerDay        sql.NullString `db:"points_per_day"`
	AcBonus             sql.NullString `db:"ac_bonus"`
	FlurryOfBlows       sql.NullString `db:"flurry_of_blows"`
	BonusSpells         sql.NullString `db:"bonus_spells"`
	PowersKnown         sql.NullString `db:"powers_known"`
	UnarmoredSpeedBonus sql.NullString `db:"unarmored_speed_bonus"`
	UnarmedDamage       sql.NullString `db:"unarmed_damage"`
	PowerLevel          sql.NullString `db:"power_level"`
	Special             sql.NullString `db:"special"`
	Slots0              sql.NullString `db:"slots_0"`
	Slots1              sql.NullString `db:"slots_1"`
	Slots2              sql.NullString `db:"slots_2"`
	Slots3              sql.NullString `db:"slots_3"`
	Slots4              sql.NullString `db:"slots_4"`
	Slots5              sql.NullString `db:"slots_5"`
	Slots6              sql.NullString `db:"slots_6"`
	Slots7              sql.NullString `db:"slots_7"`
	Slots8              sql.NullString `db:"slots_8"`
	Slots9              sql.NullString `db:"slots_9"`
	SpellsKnown0        sql.NullString `db:"spells_known_0"`
	SpellsKnown1        sql.NullString `db:"spells_known_1"`
	SpellsKnown2        sql.NullString `db:"spells_known_2"`
	SpellsKnown3        sql.NullString `db:"spells_known_3"`
	SpellsKnown4        sql.NullString `db:"spells_known_4"`
	SpellsKnown5        sql.NullString `db:"spells_known_5"`
	SpellsKnown6        sql.NullString `db:"spells_known_6"`
	SpellsKnown7        sql.NullString `db:"spells_known_7"`
	SpellsKnown8        sql.NullString `db:"spells_known_8"`
	SpellsKnown9        sql.NullString `db:"spells_known_9"`
	Reference           sql.NullString `db:"reference"`
}

var ClassTableTableName = "class_table"

type ClassTableRepo struct {
	*BaseRepo
	Class ClassListColumns
}

func (ClassTableRepo) FindAll(page uint64, pageSize uint64) (*sql.Rows, error) {
	sql, args, err := squirrel.Select("*").
		From(ClassTableTableName).
		Offset(page * pageSize).
		Limit(pageSize).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c ClassTableRepo) ExtractVars() map[string]interface{} {
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

func (c ClassTableRepo) Find(filter ClassTableColumns, page uint64, pageSize uint64) ClassTableRepo {
	c.QueryBuilder = squirrel.Select("*").From(ClassTableTableName)
	return c
}
