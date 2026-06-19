import useReveal from '../useReveal'

const statuses = [
  { name: 'MATCHED', dot: 'green', desc: 'Found in both internal records and provider data. Amounts agree within configured tolerance.' },
  { name: 'FLAGGED', dot: 'red', desc: 'Found in both sources but amounts diverge. Sent for manual review with variance logged.' },
  { name: 'MISSING_IN_PROVIDER', dot: 'red', desc: "Present internally but absent from the provider. Potential missed settlement or provider-side error." },
  { name: 'MISSING_IN_INTERNAL', dot: 'red', desc: "Provider reports a transaction that isn't in internal books. Requires investigation before period close." },
  { name: 'DUPLICATE', dot: 'red', desc: 'Same reference ID appears more than once. Engine deduplicates and flags extras automatically.' },
  { name: 'PENDING', dot: 'gray', desc: "Settlement lag is expected. Provider confirmed the transaction but the settlement window hasn't elapsed." },
]

export default function Statuses() {
  const ref = useReveal()
  return (
    <section className="section reveal" ref={ref}>
      <div className="section-num">03 — Reconciliation Statuses</div>
      <h2>Every outcome is named. Every outcome is traceable.</h2>
      <p className="section-lead">No transaction is left in an ambiguous state.</p>
      <table className="status-table">
        <tbody>
          {statuses.map((s) => (
            <tr key={s.name}>
              <td>
                <span className={`dot dot-${s.dot}`}></span>
                {s.name}
              </td>
              <td>{s.desc}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
