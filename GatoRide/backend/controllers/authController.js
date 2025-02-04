import User from '../models/User.js';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import sendVerificationEmail from '../utils/mailer.js';

const signup = async (req, res) => {
  const { name, username, email, password } = req.body;

  try {
    // Check if user already exists with the same email or username
    let user = await User.findOne({ $or: [{ email }, { username }] });
    if (user) {
      return res.status(400).json({ msg: 'User with this email or username already exists' });
    }

    // Create a new user
    user = new User({ name, username, email, password });

    // Generate a verification token
    const verificationToken = jwt.sign({ email }, process.env.JWT_SECRET, { expiresIn: '1h' });
    user.verificationToken = verificationToken;

    // Save the user to the database
    await user.save();

    // Send verification email
    await sendVerificationEmail(email, verificationToken);

    res.status(201).json({ msg: 'User registered successfully. Please check your email to verify your account.' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server error');
  }
};

const verifyEmail = async (req, res) => {
  try {
      const { token } = req.query;

      if (!token) {
          return res.status(400).json({ msg: "Invalid or missing token" });
      }

      // Verify the token
      const decoded = jwt.verify(token, process.env.JWT_SECRET);

      // Find the user by email
      const user = await User.findOne({ email: decoded.email });

      if (!user) {
          return res.status(404).json({ msg: "User not found" });
      }

      // Mark user as verified
      user.isVerified = true;
      await user.save();

      res.status(200).json({ msg: "Email verified successfully! You can now log in." });
  } catch (error) {
      console.error(error);
      res.status(400).json({ msg: "Invalid or expired token" });
  }
};

const login = async (req, res) => {
  const { usernameOrEmail, password } = req.body;

  try {
    // Find the user by username or email
    const user = await User.findOne({
      $or: [{ username: usernameOrEmail }, { email: usernameOrEmail }],
    });

    if (!user) {
      return res.status(400).json({ msg: 'user notfound' });
    }

    // Check if the password is correct
    const salt = await bcrypt.genSalt(10);
    const input_password = await bcrypt.hash(password, salt);
    const isMatch = await bcrypt.compare(password, user.password);

    if (!isMatch) {
      console.log()
      return res.status(400).json({ msg: 'Invalid credentials' });
    }

    // Check if the user's email is verified
    if (!user.isVerified) {
      return res.status(400).json({ msg: 'Please verify your email first' });
    }

    // Generate a JWT token
    const payload = { user: { id: user.id } };
    const token = jwt.sign(payload, process.env.JWT_SECRET, { expiresIn: '1h' });

    // Set token in HTTP-Only cookie
    res.cookie("token", token, {
      httpOnly: true,      // Prevent JavaScript access
      secure: true,        // Works only on HTTPS (set to false for local dev)
      sameSite: "Strict",  // Prevent CSRF attacks
      maxAge: 3600000,     // 1 hour expiration
    });

    res.json({ msg: "Login successful!" });
  } catch (err) {
    console.error(err.message);
    res.status(500).send("Server error");
  }
};

export default { signup, verifyEmail, login };