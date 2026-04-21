import { useState } from "react"
import styles from "./LogIn.module.css"
import { signInUser } from "../../auth.thunk"
import { useAppDispatch } from "../../../../app/hooks"

const LogIn = () => {
  const dispatch = useAppDispatch();

  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  
  const handleLogin = async () => {
    dispatch(signInUser({ email, password }))
  }
  return (
    <div className={styles.container}>
      <input type="email" placeholder="Email" className={styles.input} value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" className={styles.input} value={password} onChange={(e) => setPassword(e.target.value)} />
      <div className={styles.button} onClick={handleLogin}>
        <p>Login</p>
      </div>
    </div>
  )
}

export default LogIn