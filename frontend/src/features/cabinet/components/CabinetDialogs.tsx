import { DialogCloseButton, ModalDialog } from '../../../shared/components/ModalDialog';
import { CABINET_COPY, ROUTES, routeById } from '../content';

type RouteDialogProps = {
  open: boolean;
  draftId: string | null;
  error: string;
  onSelect: (routeId: string) => void;
  onConfirm: () => void;
  onCancel: () => void;
};

/** Диалог выбора маршрута — разметка из routeDialog дизайн-макета. */
export function RouteDialog({
  open,
  draftId,
  error,
  onSelect,
  onConfirm,
  onCancel,
}: RouteDialogProps) {
  return (
    <ModalDialog
      className="dialog route-dialog"
      open={open}
      onClose={onCancel}
      titleId="route-dialog-title"
      descriptionId="route-dialog-description"
    >
      <DialogCloseButton label={CABINET_COPY.actions.close} onClick={onCancel} />
      <h2 id="route-dialog-title">{CABINET_COPY.routeDialog.title}</h2>
      <p className="dialog-description" id="route-dialog-description">
        {CABINET_COPY.routeDialog.description}
      </p>
      <form
        noValidate
        onSubmit={(event) => {
          event.preventDefault();
          onConfirm();
        }}
      >
        <fieldset className="route-options">
          <legend>Дистанция</legend>
          {ROUTES.map((route) => (
            // eslint-disable-next-line jsx-a11y/label-has-associated-control -- неявный label из дизайн-макета: input внутри label, доступный текст — в span.route-option-copy (вложенность сильнее глубины поиска правила)
            <label className="route-option" key={route.id}>
              <input
                type="radio"
                name="route"
                value={route.id}
                checked={draftId === route.id}
                onChange={() => onSelect(route.id)}
              />
              <span className="route-option-copy">
                <strong>{route.distance}</strong>
                <span className="route-option-title">{route.title}</span>
                <span className="route-option-description">{route.description}</span>
              </span>
            </label>
          ))}
        </fieldset>
        {error ? (
          <p className="field-error" role="alert">
            {error}
          </p>
        ) : null}
        <div className="dialog-actions">
          <button className="primary-button" type="submit">
            {CABINET_COPY.actions.saveRoute}
          </button>
          <button className="secondary-button" type="button" onClick={onCancel}>
            {CABINET_COPY.actions.cancel}
          </button>
        </div>
      </form>
    </ModalDialog>
  );
}

type RemoveRouteDialogProps = {
  open: boolean;
  onConfirm: () => void;
  onCancel: () => void;
};

/** Подтверждение отказа от выбранного маршрута. */
export function RemoveRouteDialog({ open, onConfirm, onCancel }: RemoveRouteDialogProps) {
  return (
    <ModalDialog
      className="dialog"
      open={open}
      onClose={onCancel}
      titleId="remove-title"
      descriptionId="remove-description"
    >
      <DialogCloseButton label={CABINET_COPY.actions.close} onClick={onCancel} />
      <h2 id="remove-title">{CABINET_COPY.routeDialog.removeTitle}</h2>
      <p className="dialog-description" id="remove-description">
        {CABINET_COPY.routeDialog.removeBody}
      </p>
      <div className="dialog-actions">
        <button className="primary-button" type="button" onClick={onConfirm}>
          {CABINET_COPY.actions.removeConfirm}
        </button>
        <button className="secondary-button" type="button" onClick={onCancel}>
          {CABINET_COPY.actions.cancel}
        </button>
      </div>
    </ModalDialog>
  );
}

type SampleDialogProps = {
  open: boolean;
  login: string;
  email: string;
  routeId: string | null;
  onClose: () => void;
};

/** Черновик анкеты участника: печатается и подписывается перед стартом. */
export function SampleDialog({ open, login, email, routeId, onClose }: SampleDialogProps) {
  const route = routeById(routeId);

  return (
    <ModalDialog
      className="dialog"
      open={open}
      onClose={onClose}
      titleId="sample-title"
      descriptionId="sample-note"
    >
      <DialogCloseButton label={CABINET_COPY.actions.close} onClick={onClose} />
      <h2 id="sample-title">{CABINET_COPY.sample.title}</h2>
      <div className="sample-sheet">
        <p className="sample-watermark">{CABINET_COPY.sample.watermark}</p>
        <dl className="sample-values">
          <div>
            <dt>Логин</dt>
            <dd>{login}</dd>
          </div>
          <div>
            <dt>Электронная почта</dt>
            <dd>{email}</dd>
          </div>
          <div>
            <dt>{CABINET_COPY.sample.routeLabel}</dt>
            <dd>{route ? `${route.distance} · ${route.title}` : '—'}</dd>
          </div>
        </dl>
      </div>
      <p className="sample-note" id="sample-note">
        {CABINET_COPY.sample.note}
      </p>
      <div className="dialog-actions">
        <button
          className="primary-button"
          type="button"
          onClick={() => {
            window.print();
          }}
        >
          {CABINET_COPY.actions.printSample}
        </button>
        <button className="secondary-button" type="button" onClick={onClose}>
          {CABINET_COPY.actions.close}
        </button>
      </div>
    </ModalDialog>
  );
}
