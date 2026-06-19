import Nav from './components/Nav'
import Hero from './components/Hero'
import Problem from './components/Problem'
import HowItWorks from './components/HowItWorks'
import Statuses from './components/Statuses'
import Auditability from './components/Auditability'
import Closing from './components/Closing'
import Footer from './components/Footer'

export default function App() {
  return (
    <>
      <Nav />
      <Hero />
      <hr className="divider" />
      <Problem />
      <hr className="divider" />
      <HowItWorks />
      <hr className="divider" />
      <Statuses />
      <hr className="divider" />
      <Auditability />
      <Closing />
      <Footer />
    </>
  )
}
