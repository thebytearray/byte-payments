package service

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/model"
)

type EmailService interface {
	SendVerificationCode(toEmail, code string) error
	SendPaymentCompletionEmail(payment model.Payment, plan model.Plan) error
	SendUnderpaymentEmail(payment model.Payment, plan model.Plan, remainingAmount float64) error
	SendOverpaymentEmail(payment model.Payment, plan model.Plan, overpaidAmount float64) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (e *emailService) SendVerificationCode(toEmail, code string) error {
	em := email.NewEmail()
	em.From = fmt.Sprintf("%s <%s>", config.Cfg.EMAIL_FROM_NAME, config.Cfg.EMAIL_FROM_ADDR)
	em.To = []string{toEmail}
	em.Subject = "BytePayments - Email Verification Code"
	
	em.HTML = []byte(fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #007bff; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.code { font-size: 24px; font-weight: bold; color: #007bff; letter-spacing: 3px; }
				.footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>BytePayments</h1>
				</div>
				<div class="content">
					<h2>Email Verification Required</h2>
					<p>Please use the following verification code to complete your payment:</p>
					<div class="code">%s</div>
					<p>This code will expire in 10 minutes.</p>
					<p>If you didn't request this code, please ignore this email.</p>
				</div>
				<div class="footer">
					<p>© 2025 BytePayments. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, code))

	// SMTP server configuration
	auth := smtp.PlainAuth("", config.Cfg.EMAIL_USERNAME, config.Cfg.EMAIL_PASSWORD, config.Cfg.EMAIL_SMTP_HOST)
	
	return em.Send(fmt.Sprintf("%s:%d", config.Cfg.EMAIL_SMTP_HOST, config.Cfg.EMAIL_SMTP_PORT), auth)
}

func (e *emailService) SendPaymentCompletionEmail(payment model.Payment, plan model.Plan) error {
	template, err := e.loadTemplate("static/email_payment_completed.html")
	if err != nil {
		return fmt.Errorf("failed to load completion email template: %w", err)
	}

	completionDate := time.Now().Format("January 2, 2006 at 3:04 PM MST")
	
	// Safety check for plan name
	planName := plan.Name
	if planName == "" {
		planName = "Selected Plan" // fallback
	}
	
	replacements := map[string]string{
		"{{PAYMENT_ID}}":      payment.ID,
		"{{PLAN_NAME}}":       planName,
		"{{AMOUNT_PAID}}":     fmt.Sprintf("%.2f", payment.PaidAmountTRX),
		"{{USER_EMAIL}}":      payment.UserEmail,
		"{{COMPLETION_DATE}}": completionDate,
	}

	// Debug logging
	fmt.Printf("DEBUG: Payment completion email replacements: %+v\n", replacements)

	htmlContent := e.replaceTemplateVars(template, replacements)

	em := email.NewEmail()
	em.From = fmt.Sprintf("%s <%s>", config.Cfg.EMAIL_FROM_NAME, config.Cfg.EMAIL_FROM_ADDR)
	em.To = []string{payment.UserEmail}
	em.Subject = "Payment Completed - BytePayments"
	em.HTML = []byte(htmlContent)

	auth := smtp.PlainAuth("", config.Cfg.EMAIL_USERNAME, config.Cfg.EMAIL_PASSWORD, config.Cfg.EMAIL_SMTP_HOST)
	return em.Send(fmt.Sprintf("%s:%d", config.Cfg.EMAIL_SMTP_HOST, config.Cfg.EMAIL_SMTP_PORT), auth)
}

func (e *emailService) SendUnderpaymentEmail(payment model.Payment, plan model.Plan, remainingAmount float64) error {
	template, err := e.loadTemplate("static/email_payment_underpaid.html")
	if err != nil {
		return fmt.Errorf("failed to load underpayment email template: %w", err)
	}

	expiryTime := payment.CreatedAt.Add(15 * time.Minute).Format("January 2, 2006 at 3:04 PM MST")
	
	// Safety checks
	planName := plan.Name
	if planName == "" {
		planName = "Selected Plan" // fallback
	}
	
	walletAddress := payment.Wallet.WalletAddress
	if walletAddress == "" {
		walletAddress = "Wallet address not available" // fallback
	}
	
	replacements := map[string]string{
		"{{PAYMENT_ID}}":       payment.ID,
		"{{PLAN_NAME}}":        planName,
		"{{REQUIRED_AMOUNT}}":  fmt.Sprintf("%.2f", payment.AmountTRX),
		"{{PAID_AMOUNT}}":      fmt.Sprintf("%.2f", payment.PaidAmountTRX),
		"{{REMAINING_AMOUNT}}": fmt.Sprintf("%.2f", remainingAmount),
		"{{WALLET_ADDRESS}}":   walletAddress,
		"{{EXPIRY_TIME}}":      expiryTime,
	}

	htmlContent := e.replaceTemplateVars(template, replacements)

	em := email.NewEmail()
	em.From = fmt.Sprintf("%s <%s>", config.Cfg.EMAIL_FROM_NAME, config.Cfg.EMAIL_FROM_ADDR)
	em.To = []string{payment.UserEmail}
	em.Subject = "Payment Incomplete - Action Required - BytePayments"
	em.HTML = []byte(htmlContent)

	auth := smtp.PlainAuth("", config.Cfg.EMAIL_USERNAME, config.Cfg.EMAIL_PASSWORD, config.Cfg.EMAIL_SMTP_HOST)
	return em.Send(fmt.Sprintf("%s:%d", config.Cfg.EMAIL_SMTP_HOST, config.Cfg.EMAIL_SMTP_PORT), auth)
}

func (e *emailService) SendOverpaymentEmail(payment model.Payment, plan model.Plan, overpaidAmount float64) error {
	template, err := e.loadTemplate("static/email_payment_overpaid.html")
	if err != nil {
		return fmt.Errorf("failed to load overpayment email template: %w", err)
	}

	completionDate := time.Now().Format("January 2, 2006 at 3:04 PM MST")
	
	// Safety check for plan name
	planName := plan.Name
	if planName == "" {
		planName = "Selected Plan" // fallback
	}
	
	replacements := map[string]string{
		"{{PAYMENT_ID}}":       payment.ID,
		"{{PLAN_NAME}}":        planName,
		"{{REQUIRED_AMOUNT}}":  fmt.Sprintf("%.2f", payment.AmountTRX),
		"{{AMOUNT_PAID}}":      fmt.Sprintf("%.2f", payment.PaidAmountTRX),
		"{{USER_EMAIL}}":       payment.UserEmail,
		"{{COMPLETION_DATE}}":  completionDate,
		"{{OVERPAID_AMOUNT}}":  fmt.Sprintf("%.2f", overpaidAmount),
	}

	htmlContent := e.replaceTemplateVars(template, replacements)

	em := email.NewEmail()
	em.From = fmt.Sprintf("%s <%s>", config.Cfg.EMAIL_FROM_NAME, config.Cfg.EMAIL_FROM_ADDR)
	em.To = []string{payment.UserEmail}
	em.Subject = "Payment Completed (Overpaid) - BytePayments"
	em.HTML = []byte(htmlContent)

	auth := smtp.PlainAuth("", config.Cfg.EMAIL_USERNAME, config.Cfg.EMAIL_PASSWORD, config.Cfg.EMAIL_SMTP_HOST)
	return em.Send(fmt.Sprintf("%s:%d", config.Cfg.EMAIL_SMTP_HOST, config.Cfg.EMAIL_SMTP_PORT), auth)
}

func (e *emailService) loadTemplate(templatePath string) (string, error) {
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (e *emailService) replaceTemplateVars(template string, replacements map[string]string) string {
	result := template
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
} 