import useReveal from '../useReveal'

export default function Closing() {
  const ref = useReveal()
  return (
    <div className="closing reveal" ref={ref}>
      <div className="closing-inner">
        <h2>
          This is not
          <br />
          a CRUD API.
          <br />
          It is <em>financial infrastructure.</em>
        </h2>
        <p>
          smart-reconcile is actively being built. The engine, normalizer, and provider
          interface are under development.
        </p>
      </div>
    </div>
  )
}
