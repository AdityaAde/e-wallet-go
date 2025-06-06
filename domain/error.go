package domain

import "errors"

var ErrAuthFailed = errors.New("Error auth failed")
var ErrUsernameTaken = errors.New("Username already taken")
var ErrOtpInvalid = errors.New("Otp invalid")
var ErrAccountNotFound = errors.New("Account not found")
var ErrInquiryNotFound = errors.New("Inquiry not found")
var ErrInsufficientBalance = errors.New("Insufficient balance")
