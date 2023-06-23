package structs

type InferenceStruct struct {
	UniqueModelName string                 `json:"model_name"`
	Params          map[string]interface{} `json:"params"`
}

type GenericStruct struct {
	JobTypeID string                 `json:"job_type_id"`
	Params    map[string]interface{} `json:"params"`
}
