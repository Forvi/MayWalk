/**
 * Знак «Майской прогулки»: восемь лепестков-эллипсов вокруг сердцевины.
 * Цвет наследуется через currentColor, подложка — через var(--canvas),
 * поэтому знак одинаково работает в паспорте и на сцене.
 */
export function BrandMark() {
  return (
    <svg className="brand-mark" viewBox="0 0 44 44" aria-hidden="true">
      <g fill="currentColor">
        {[0, 45, 90, 135, 180, 225, 270, 315].map((angle) => (
          <ellipse
            key={angle}
            cx="22"
            cy="10.5"
            rx="5.8"
            ry="10.5"
            transform={`rotate(${angle} 22 22)`}
          />
        ))}
        <circle cx="22" cy="22" r="10" />
      </g>
      <circle cx="22" cy="22" r="4.2" fill="var(--canvas,#F2F1E7)" />
    </svg>
  );
}

type FlowerMarkProps = {
  /** Класс из дизайн-макета: lead-mark (лендинг) или route-flower (кабинет). */
  className?: string;
  /** Сторона в пикселях: 105 — лендинг, 300 — витрина маршрута. */
  size?: number;
};

/**
 * Крупный цветок. Класс выбирает место применения:
 * на лендинге это .lead-mark (absolute, цвет наследуется темой),
 * в кабинете — .route-flower (фон витрины маршрута, красится родителем).
 */
export function FlowerMark({ className = 'route-flower', size = 300 }: FlowerMarkProps) {
  return (
    <svg
      className={className}
      width={size}
      height={size}
      viewBox="0 0 300 300"
      aria-hidden="true"
    >
      <g fill="currentColor">
        {[0, 36, 72, 108, 144, 180, 216, 252, 288, 324].map((angle) => (
          <ellipse
            key={angle}
            cx="150"
            cy="73"
            rx="33"
            ry="69"
            transform={`rotate(${angle} 150 150)`}
          />
        ))}
        <circle cx="150" cy="150" r="78" />
      </g>
    </svg>
  );
}
