import User from '../models/User.js';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import sendVerificationEmail from '../utils/mailer.js';

const signup = async (req, res) => {
  const { name, username, email, password } = req.body;

  try {
    // Check if user already exists with the same email or username
    let user = await User.findOne({ $or: [{ email }, { username }] });
    // if (user) {
    //   return res.status(400).json({ msg: 'User with this email or username already exists' });
    // }

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
  const { token } = req.query;

  try {
    // Verify the token
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    const user = await User.findOne({ email: decoded.email });

    if (!user) {
      return res.status(400).json({ msg: 'Invalid token' });
    }

    // Mark the user as verified
    user.isVerified = true;
    user.verificationToken = undefined;
    await user.save();

    res.status(200).json({ msg: 'Email verified successfully' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server error');
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
      return res.status(400).json({ msg: 'Invalid credentials' });
    }

    // Check if the password is correct
    const isMatch = await bcrypt.compare(password, user.password);
    if (!isMatch) {
      return res.status(400).json({ msg: 'Invalid credentials' });
    }

    // Check if the user's email is verified
    if (!user.isVerified) {
      return res.status(400).json({ msg: 'Please verify your email first' });
    }

    // Generate a JWT token
    const payload = {
      user: {
        id: user.id,
      },
    };

    jwt.sign(payload, process.env.JWT_SECRET, { expiresIn: '1h' }, (err, token) => {
      if (err) throw err;
      res.json({ token });
    });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server error');
  }
};

export default { signup, verifyEmail, login };