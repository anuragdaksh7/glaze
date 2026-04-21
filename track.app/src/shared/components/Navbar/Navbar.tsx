import styles from "./Navbar.module.css"

type Props = {
  children: React.ReactNode
}

const Navbar = ({ children }: Props) => {
  return (
    <nav className={styles.nav}>{children}</nav>
  )
}

export default Navbar