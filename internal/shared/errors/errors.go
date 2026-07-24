package errors

import "errors"

// employees error type
var ErrEmployeeNotFound error = errors.New("employee not found")

var ErrEmployeeMaserNotFound error = errors.New("master not found")

// branch error type
var ErrBranchNotFound error = errors.New("branch not found")

// department error type
var ErrDepartmentNotFound error = errors.New("department not found")

var ErrCategoryNotFound error = errors.New("category not found")

var ErrPositionNotFound error = errors.New("position not found")

var ErrEmployeeOrCategoryNotFound error = errors.New("employee or category not found")

var ErrCategoryHasPosition error = errors.New("category has position")

var ErrMaterialUnitCategoryHasMaterialUnit error = errors.New("material unit category has material unit")

var ErrTypeNotFound error = errors.New("type not found")

var ErrUserNotFound error = errors.New("user not found")

var ErrManufacturerNotFound error = errors.New("manufacturer not found")

var ErrCountryNotFound error = errors.New("country not found")

var ErrCompanyNotFound error = errors.New("company not found")

var ErrPurchasingOrganizationNotFound error = errors.New("purchasing organization not found")

var ErrPurchasingGroupNotFound error = errors.New("purchasing group not found")

var ErrPurchasingGroupMemberNotFound error = errors.New("purchasing group member not found")

var ErrLanguageNotFound error = errors.New("language not found")

var ErrCurrencyNotFound error = errors.New("currency not found")

var ErrPaymentTermNotFound error = errors.New("payment term not found")

var ErrPaymentMethodNotFound error = errors.New("payment method not found")

var ErrStatusNotFound error = errors.New("status not found")

var ErrBrandNotFound error = errors.New("brand not found")

var ErrMasterUserNotFound error = errors.New("master not found")

var ErrMutationError error = errors.New("mutation error")

var ErrBlankPassword error = errors.New("password is blank")

var ErrUserLocked error = errors.New("user is locked")

var ErrInvalidPassword error = errors.New("invalid password")

var ErrTooManyRequest error = errors.New("too many request")

var ErrInvalidSession error = errors.New("invalid session")

var ErrDuplicateError error = errors.New("duplicate error")

var ErrDuplicateUsernameError error = errors.New("duplicate username")

var ErrDuplicateLegalError error = errors.New("duplicate legal")

var ErrDuplicateContactError error = errors.New("duplicate contact")

var ErrDuplicateAttendanceError error = errors.New("duplicate attendance")

var ErrTodayIsWeekend error = errors.New("today is weekend")

var ErrAttendanceNotFound error = errors.New("attendance not found")

var ErrWarehouseNotFound error = errors.New("warehouse not found")

var ErrVendorNotFound error = errors.New("vendor not found")

var ErrVendorContactNotFound error = errors.New("vendor contact not found")

var ErrDuplicateVendorContact error = errors.New("duplicate vendor contact")

var ErrMaterialTypeNotFound error = errors.New("material type not found")

var ErrMaterialUnitNotFound error = errors.New("material unit not found")

var ErrBaseUnitNotFound error = errors.New("base unit not found")

var ErrWeightUnitNotFound error = errors.New("weight unit not found")

var ErrVolumeUnitNotFound error = errors.New("volume unit not found")

var ErrMaterialGroupNotFound error = errors.New("material group not found")

var ErrStorageLocationNotFound error = errors.New("storage location not found")

var ErrDuplicateStorageLocationCode error = errors.New("storage location not found")

var ErrIndustrySectorNotFound error = errors.New("industry sector not found")

var ErrPositionCategoryHasActivePosition error = errors.New("position category has active positions")

var ErrBranchHasActiveEmployees error = errors.New("branch has active employees")

var ErrPositionHasActiveEmployees error = errors.New("position has active employees")

var ErrParentPositionHasActiveChild error = errors.New("parent position has active child")

var ErrDepartmentHasActivePosition error = errors.New("department has active position")

var ErrForeignKeyError error = errors.New("some id missing or not found")

var ErrMaterialMasterNotFound error = errors.New("material master not found")
