package model

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type Account struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account Account `json:"account"`
}

type Tokenclaim struct {
	Exp               int            `json:"exp"`
	Iat               int            `json:"iat"`
	AuthTime          int            `json:"auth_time"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               interface{}    `json:"aud"`
	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	SessionState      string         `json:"session_state"`
	Acr               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	EmailVerified     bool           `json:"email_verified"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
}

func (t *Tokenclaim) AudAsSlice() []string {
	switch t.Aud.(type) {
	case string:
		return []string{t.Aud.(string)}
	case []interface{}:
		auds,ok := t.Aud.([]interface{})
		if !ok {
			return []string{}
		}
		result := []string{}
		for _, aud := range auds {
			if sAud,ok := aud.(string); ok {
				result = append(result, sAud)
			}
		}
		return result
	default:
		return []string{}
	}
}
