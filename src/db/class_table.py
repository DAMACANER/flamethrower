from sqlalchemy import Column, Integer, String
from src.definitions import Base
from sqlalchemy.ext.asyncio import AsyncSession
from typing import Sequence
from sqlalchemy import select


class ClassTable(Base):
    """
    Example:

    An example instance of ClassTable in JSON format would look like this:

    ```json
        {
            "flurry_of_blows": "None",
            "slots_1": "None",
            "slots_8": "None",
            "spells_known_5": "None",
            "name": "Barbarian",
            "bonus_spells": "None",
            "slots_2": "None",
            "slots_9": "None",
            "spells_known_7": "None",
            "fort_save": "+7",
            "powers_known": "None",
            "slots_3": "None",
            "spells_known_0": "None",
            "spells_known_8": "None",
            "ref_save": "+3",
            "unarmored_speed_bonus": "None",
            "slots_4": "None",
            "spells_known_1": "None",
            "spells_known_9": "None",
            "will_save": "+3",
            "unarmed_damage": "None",
            "slots_5": "None",
            "spells_known_2": "None",
            "reference": "SRD 3.5 ClassesI",
            "caster_level": "None",
            "power_level": "None",
            "slots_6": "None",
            "spells_known_3": "None",
            "level": "11",
            "points_per_day": "None",
            "special": "Greater rage",
            "slots_7": "None",
            "spells_known_4": "None",
            "id": 11,
            "base_attack_bonus": "+11/+6/+1",
            "ac_bonus": "None",
            "slots_0": "None",
            "spells_known_6": "None"
        },
    ```
    """

    __tablename__ = "class_table"
    id = Column("id", Integer, primary_key=True)
    name = Column("name", String)
    level = Column("level", String)
    base_attack_bonus = Column("base_attack_bonus", String)
    fort_save = Column("fort_save", String)
    ref_save = Column("ref_save", String)
    will_save = Column("will_save", String)
    caster_level = Column("caster_level", String)
    points_per_day = Column("points_per_day", String)
    ac_bonus = Column("ac_bonus", String)
    flurry_of_blows = Column("flurry_of_blows", String)
    bonus_spells = Column("bonus_spells", String)
    powers_known = Column("powers_known", String)
    unarmored_speed_bonus = Column("unarmored_speed_bonus", String)
    unarmed_damage = Column("unarmed_damage", String)
    power_level = Column("power_level", String)
    special = Column("special", String)
    slots_0 = Column("slots_0", String)
    slots_1 = Column("slots_1", String)
    slots_2 = Column("slots_2", String)
    slots_3 = Column("slots_3", String)
    slots_4 = Column("slots_4", String)
    slots_5 = Column("slots_5", String)
    slots_6 = Column("slots_6", String)
    slots_7 = Column("slots_7", String)
    slots_8 = Column("slots_8", String)
    slots_9 = Column("slots_9", String)
    spells_known_0 = Column("spells_known_0", String)
    spells_known_1 = Column("spells_known_1", String)
    spells_known_2 = Column("spells_known_2", String)
    spells_known_3 = Column("spells_known_3", String)
    spells_known_4 = Column("spells_known_4", String)
    spells_known_5 = Column("spells_known_5", String)
    spells_known_6 = Column("spells_known_6", String)
    spells_known_7 = Column("spells_known_7", String)
    spells_known_8 = Column("spells_known_8", String)
    spells_known_9 = Column("spells_known_9", String)
    reference = Column("reference", String)


class ClassTableRepository:
    def __init__(self, db: AsyncSession):
        self.db = db

    async def get_all(self, page: int = 1, page_size: int = 10) -> Sequence[ClassTable]:
        query = select(ClassTable).offset(page_size * page).limit(page_size)
        return (await self.db.execute(query)).scalars().all()
