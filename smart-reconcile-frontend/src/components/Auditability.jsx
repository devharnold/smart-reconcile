import useReveal from '../useReveal'

const questions = [
  {
    q: 'Why was this transaction marked as reconciled?',
    a: <>Matched by reference ID and amount within tolerance. <strong>Rule ID and run ID persisted.</strong></>,
  },
  {
    q: 'What was the raw payload from the provider at the time?',
    a: <>Stored verbatim in <strong>raw_transactions</strong> before any transformation. Always retrievable.</>,
  },
  {
    q: 'Who flagged this for manual review and when?',
    a: <>Engine rule: <strong>amount_variance &gt; configured threshold.</strong> Timestamp and run ID logged.</>,
  },
  {
    q: 'What was the variance?',
    a: <>Exact decimal value stored per reconciliation record. <strong>Not rounded. Not approximated.</strong></>,
  },
]

export default function Auditability() {
  const ref = useReveal()
  return (
    <section className="section reveal" ref={ref}>
      <div className="section-num">04 — Auditability</div>
      <h2>Without answers, reconciliation reports are just numbers.</h2>
      <p className="section-lead">
        In financial systems, you must be able to answer these questions. With them, reports
        become a legal record.
      </p>
      <div className="audit-questions">
        {questions.map((item) => (
          <div className="audit-q" key={item.q}>
            <span className="q">{item.q}</span>
            <span className="a">{item.a}</span>
          </div>
        ))}
      </div>
    </section>
  )
}
