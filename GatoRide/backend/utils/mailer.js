import nodemailer from 'nodemailer';

const transporter = nodemailer.createTransport({
  host: "smtp.gmail.com",
  port: 465,
  secure: true,
  auth: {
    user: "shahidshareef4457@gmail.com",
    pass: "ldvqpjhklisopenq",
  },
});

const sendVerificationEmail = async (email, token) => {
  console.log("SMTP USER:", `"${process.env.EMAIL_USER}"`);
  console.log("SMTP PASS:", `"${process.env.EMAIL_PASS}"`);

  const mailOptions = {
    from: process.env.EMAIL_USER,
    to: email,
    subject: "Email Verification",
    text: `Verify your email by clicking the link: http://localhost:3000/verify-email?token=${token}`,
  };

  try {
    await transporter.sendMail(mailOptions);
    console.log("✅ Email sent successfully");
  } catch (error) {
    console.error("❌ Email sending failed:", error);
  }
};

export default sendVerificationEmail;
