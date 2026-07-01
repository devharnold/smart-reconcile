import useReveal from '../useReveal'

const stages = [
  {
    title: 'Merchant / internal system',
    tag: 'Source',
    desc: 'The transaction originates here — a sale, a payout, an internal ledger entry that needs to be reconciled against an external provider.',
  },
  {
    title: 'Provider integration layer',
    tag: 'M-Pesa · Stripe · PayPal · Bank',
    desc: 'Each provider is integrated independently behind a shared interface. Provider-specific auth, pagination, and payload quirks live here and nowhere else.',
  },
  {
    title: 'Polling & scheduler engine',
    tag: 'Cloud Scheduler → Pub/Sub → scheduler-worker',
    desc: 'Cloud Scheduler triggers scheduler jobs on a fixed cadence. Each job is published to Pub/Sub and picked up by a scheduler-worker, decoupling the act of fetching from the act of processing.',
  },
  {
    title: 'Raw transaction storage',
    tag: 'JSONB, untouched',
    desc: 'The original payload is stored exactly as received, before any transformation. This is the immutable record that everything downstream can be checked against.',
  },
  {
    title: 'Normalization pipeline',
    tag: 'normalize-worker',
    desc: 'A dedicated worker maps every provider payload onto a single shared transaction model — consistent fields, consistent units, consistent timestamps.',
  },
  {
    title: 'Reconciliation engine',
    tag: 'Reference · tolerance · fuzzy matching',
    desc: 'Normalized transactions are matched against internal records using a layered strategy: exact reference match first, then amount tolerance, then fuzzy matching for edge cases.',
  },
  {
    title: 'Mismatch detection & result storage',
    tag: 'Status + reasoning persisted',
    desc: 'Every match attempt produces a result — matched, flagged, or missing — along with the reasoning behind it. Nothing is discarded.',
  },
  {
    title: 'Reports · dashboard · alerts',
    tag: 'Output',
    desc: 'Reconciliation results surface as reports, a live dashboard, and alerts for anything that needs human attention.',
  },
]

export default function HowItWorks() {
  const ref = useReveal()
  return (
    <section className="section reveal" ref={ref}>
      <div className="section-num">02 — How It Works</div>
      <h2>Eight stages. Each one decoupled from the next.</h2>
      <p className="section-lead">
        Every stage is handled by a dedicated worker, connected via Pub/Sub rather than direct
        calls. Data moves through the pipeline one stage at a time, and each stage only knows
        about the one before it.
      </p>
      <div className="pipeline">
        {stages.map((s, i) => (
          <div className="pipeline-stage" key={s.title}>
            <div className="stage-n">{String(i + 1).padStart(2, '0')}</div>
            <div className="stage-content">
              <h3>
                {s.title}
                <span className="stage-tag">{s.tag}</span>
              </h3>
              <p>{s.desc}</p>
            </div>
          </div>
        ))}
      </div>
      <div className="pipeline-note">
        <ul>
          <li>Stages can fail independently and retry without cascading into the rest of the pipeline.</li>
          <li>Each worker scales horizontally on its own, without affecting any other stage.</li>
          <li>There is a natural audit boundary between raw ingestion and processed output — nothing is normalized until the original is safely stored.</li>
        </ul>
      </div>
    </section>
  )
}
