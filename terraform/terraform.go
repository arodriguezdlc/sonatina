package terraform

type Terraform interface {
	Init()
	Plan()
	Apply()
	ApplyPlan()
	Destroy()
}
