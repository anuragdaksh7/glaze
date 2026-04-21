import styles from "./Logout.module.css"

const Logout = () => {
  const handleLogout = () => {
    localStorage.removeItem("user")
    window.location.href = "/login"
  }

  return (
    <div onClick={handleLogout} className={styles.logout_button}>
      <p>Logout</p>
    </div>
  )
}

export default Logout