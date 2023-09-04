/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var microdevlog = logf.Log.WithName("microdev-resource")

func (r *MicroDev) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-test-606-hdu-io-v1-microdev,mutating=true,failurePolicy=fail,sideEffects=None,groups=test.606.hdu.io,resources=microdevs,verbs=create;update,versions=v1,name=mmicrodev.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MicroDev{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MicroDev) Default() {
	microdevlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.QueryInterval == nil {
		r.Spec.QueryInterval = new(int32)
		*r.Spec.QueryInterval = 30
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-test-606-hdu-io-v1-microdev,mutating=false,failurePolicy=fail,sideEffects=None,groups=test.606.hdu.io,resources=microdevs,verbs=create;update,versions=v1,name=vmicrodev.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MicroDev{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MicroDev) ValidateCreate() (admission.Warnings, error) {
	microdevlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.validateMicroDev()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MicroDev) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	microdevlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, r.validateMicroDev()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MicroDev) ValidateDelete() (admission.Warnings, error) {
	microdevlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *MicroDev) validateMicroDev() error {
	var allErrs field.ErrorList
	if err := r.validateMicroDevName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateMicroDevSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "test.606.hdu.io", Kind: "MicroDev"},
		r.Name, allErrs)
}

func (r *MicroDev) validateMicroDevName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 character like all Kubernetes objects
		// (which must fit in a DNS subdomain). The MicroDev controller appends
		// a 11-character suffix to the MicroDev (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore MicroDev
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 52 characters")
	}
	return nil
}

func (r *MicroDev) validateMicroDevSpec() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	return validateQueryUrlFormat(r.Spec.QueryUrl)
}

// 验证url格式是否正确
func validateQueryUrlFormat(queryUrl string) *field.Error {
	return nil
}
