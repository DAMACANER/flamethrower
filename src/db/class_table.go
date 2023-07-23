package db

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
)

type ClassTable struct {
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

var ClassTableName = "class_table"

type ClassTableRepo struct{}

func (ClassTableRepo) FindAll() (*sql.Rows, error) {
	sql, args, err := squirrel.Select("*").From(ClassTableName).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (ClassTableRepo) Find(filter ClassTable) (*sql.Rows, error) {
	q := squirrel.Select("*").From(ClassTableName)

	v := reflect.ValueOf(filter)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("db")
		switch field.Interface().(type) {
		case sql.NullString:
			ns := field.Interface().(sql.NullString)
			if ns.Valid {
				q = q.Where(squirrel.Eq{tag: ns.String})
			}
		case string:
			if field.String() != "" {
				q = q.Where(squirrel.Eq{tag: field.String()})
			}
		case int:
			if field.Int() != 0 {
				q = q.Where(squirrel.Eq{tag: fmt.Sprintf("%d", field.Int())})
			}
		}
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
