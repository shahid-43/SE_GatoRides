import nodemailer from 'nodemailer';

const transporter = nodemailer.createTransport({
  service: 'gmail',
  auth: {
    user: process.env.EMAIL_USER,
    pass: process.env.EMAIL_PASS,
  },
});

const sendVerificationEmail = async (email, token) => {
    console.log("SMTP USER:", process.env.EMAIL_USER);
    console.log("SMTP PASS:", process.env.EMAIL_PASS);
  const mailOptions = {
    from: process.env.EMAIL_USER,
    to: email,
    subject: 'Email Verification',
    text: `Please verify your email by clicking the following link: http://localhost:3000/verify-email?token=${token}`,
  };

  await transporter.sendMail(mailOptions);
};

export default sendVerificationEmail;