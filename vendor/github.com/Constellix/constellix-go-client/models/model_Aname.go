package models

type AnameAttributes struct {
	Name                string            `json:"name,omitempty"`
	TTL                 int               `json:"ttl"`
	RecordOption        string            `json:"recordOption,omitempty"`
	NoAnswer            bool              `json:"noAnswer,"`
	SkipLookup          bool              `json:"skipLookup,"`
	Note                string            `json:"note,omitempty"`
	GtdRegion           int               `json:"gtdRegion,omitempty"`
	Type                string            `json:"type,omitempty"`
	ContactIDs          []int             `json:"contactIds,omitempty"`
	RoundRobin          []interface{}     `json:"roundRobin,omitempty"`
	Pools               []int             `json:"pools,omitempty"`
	GeoLocation         *GeolocationANAME `json:"geolocation,omitempty"`
	RecordFailoverAname *RCDFAname        `json:"recordFailover,omitempty"`
}

type GeolocationANAME struct {
	GeoIpUserRegion []int `json:"geoipUserRegion,omitempty"`
	Drop            bool  `json:"drop,"`
	GeoIpProximity  int   `json:"geoipProximity,omitempty"`
	GeoIpFailOver   bool  `json:"geoipFailover,"`
}

type AnameRoundRobin struct {
	Value       string `json:"value,omitempty"`
	DisableFlag string `json:"disableFlag,omitempty"`
}

type RCDFAname struct {
	Values            []interface{} `json:"values,omitempty"`
	FailoverTypeAname int           `json:"failoverType,omitempty"`
	DisableFlagAname  bool          `json:"disabled,"`
}

type ValuesAname struct {
	Value         string `json:"value,omitempty"`
	SortOrderRCDF int    `json:"sortOrder,omitempty"`
	DisableFlag   bool   `json:"disableFlag,"`
	CheckID       int    `json:"checkid,omitempty"`
}
