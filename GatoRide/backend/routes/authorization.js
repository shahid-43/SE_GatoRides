import express from 'express';
import authController from '../controllers/authController.js';

const { signup, verifyEmail, login } = authController;

const router = express.Router();

router.post('/signup', signup);
router.get('/verify-email', verifyEmail);
router.post('/login', login);

export default router;