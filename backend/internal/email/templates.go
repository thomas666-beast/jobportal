package email

import "fmt"

const layout = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>%s</title></head>
<body style="margin:0;padding:0;background:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="padding:32px 16px;">
    <tr><td align="center">
      <table role="presentation" width="100%%" style="max-width:520px;background:#ffffff;border-radius:12px;border:1px solid #e2e8f0;padding:32px;">
        <tr><td style="font-size:20px;font-weight:600;color:#0f172a;padding-bottom:8px;">%s</td></tr>
        <tr><td style="font-size:24px;font-weight:600;color:#0f172a;padding:16px 0 8px;">%s</td></tr>
        <tr><td style="font-size:15px;color:#475569;line-height:1.6;padding:8px 0 24px;">%s</td></tr>
        <tr><td style="padding:8px 0;">
          <a href="%s" style="display:inline-block;background:#4f46e5;color:#ffffff;text-decoration:none;padding:12px 24px;border-radius:8px;font-weight:500;font-size:15px;">%s</a>
        </td></tr>
        <tr><td style="font-size:13px;color:#94a3b8;padding:24px 0 0;line-height:1.5;">%s</td></tr>
      </table>
      <p style="font-size:12px;color:#94a3b8;margin-top:24px;">This email was sent by %s. If you didn't request it, you can safely ignore it.</p>
    </td></tr>
  </table>
</body></html>`

type RenderedEmail struct {
	Subject string
	HTML    string
	Text    string
}

func RenderVerificationEmail(appName, baseURL, token string) RenderedEmail {
	link := fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)

	html := fmt.Sprintf(layout,
		"Verify your email",                                                                                // 1: title
		appName,                                                                                            // 2: header brand
		"Verify your email",                                                                                // 3: heading
		fmt.Sprintf("Welcome to %s. Please verify your email address to unlock applying to jobs and posting listings.", appName), // 4: body
		link,                                                                                               // 5: href
		"Verify email",                                                                                     // 6: button
		"This link expires in 24 hours.",                                                                   // 7: footnote
		appName,                                                                                            // 8: footer brand
	)

	text := fmt.Sprintf(`Welcome to %s.

Please verify your email address by clicking the link below:
%s

This link expires in 24 hours.

If you didn't create an account, ignore this email.`, appName, link)

	return RenderedEmail{
		Subject: fmt.Sprintf("Verify your email — %s", appName),
		HTML:    html,
		Text:    text,
	}
}

func RenderPasswordResetEmail(appName, baseURL, token string) RenderedEmail {
	link := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	html := fmt.Sprintf(layout,
		"Reset your password",                                                                                       // 1: title
		appName,                                                                                                     // 2: header brand
		"Reset your password",                                                                                       // 3: heading
		fmt.Sprintf("We received a request to reset your %s password. Click below to set a new one.", appName),      // 4: body
		link,                                                                                                        // 5: href
		"Reset password",                                                                                            // 6: button
		"This link expires in 1 hour. If you didn't request it, ignore this email — your password won't change.",     // 7: footnote
		appName,                                                                                                     // 8: footer brand
	)

	text := fmt.Sprintf(`Someone (hopefully you) requested a password reset for your %s account.

Click below to set a new password:
%s

This link expires in 1 hour.

If you didn't request it, ignore this email — your password won't change.`, appName, link)

	return RenderedEmail{
		Subject: fmt.Sprintf("Reset your password — %s", appName),
		HTML:    html,
		Text:    text,
	}
}
