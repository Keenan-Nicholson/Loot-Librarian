package main

type DamageType struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Damage struct {
	DamageDice string     `json:"damage_dice"`
	DamageType DamageType `json:"damage_type"`
}

type Cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type Range struct {
	Normal int `json:"normal"`
}

type ArmorType struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

type Equipment struct {
	WeaponCategory     string `json:"weapon_category"`
	WeaponRange        string `json:"weapon_range"`
	Damage             Damage `json:"damage"`
	WeaponRangeDetails Range  `json:"range"` // Changed field name to avoid conflict with WeaponRange

	ArmorCategory       string     `json:"armor_category"`
	ArmorClass          ArmorClass `json:"armor_class"`
	StrMinimum          int        `json:"str_minimum"`
	StealthDisadvantage bool       `json:"stealth_disadvantage"`
	ArmorWeight         int        `json:"weight"` // Changed field name to avoid conflict with Armor Weight

	Cost Cost `json:"cost"`
}
