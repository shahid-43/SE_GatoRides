import mongoose from 'mongoose';
import bcrypt from 'bcryptjs';

const userSchema = new mongoose.Schema({
    name: {
        type: String,
        required: [true, "Please provide a name"],
    },
    email: {
        type: String,
        required: [true, "Please provide an email"],
        unique: true,
    },
    username: {
        type: String,
        required: [true, "Please provide a username"],
        unique: true,
        immutable: true, // Ensures that username cannot be changed
    },
    password: {
        type: String,
        required: [true, "Please provide a password"],
    },
    isVerified: {
        type: Boolean,
        default: false,
      },
      verificationToken: {
        type: String,
      },
    });

    userSchema.pre("save", async function (next) {
        // if (this.isModified("username")) {
        //     return next(new Error("Username cannot be modified"));
        // }
        if (!this.isModified("password")) {
          next();
        }
        const salt = await bcrypt.genSalt(10);
        this.password = await bcrypt.hash(this.password, salt);
      } );

// module.exports = mongoose.models.User || mongoose.model("User", userSchema);
// User.js
export default mongoose.model('User', userSchema);
