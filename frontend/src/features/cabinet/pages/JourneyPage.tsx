import { useState } from 'react';

import { useAuth } from '../../../features/auth/AuthContext';
import { FlowerMark } from '../../../shared/components/BrandMark';
import { CABINET_COPY, routeById } from '../content';
import { RemoveRouteDialog, RouteDialog, SampleDialog } from '../components/CabinetDialogs';

type DialogKind = 'route' | 'remove' | 'sample' | null;

/**
 * «Моя прогулка»: витрина выбранного маршрута, смена маршрута, черновик анкеты.
 * Логика состояний перенесена из src/v4/cabinet-state.js дизайн-макета.
 */
export function JourneyPage() {
  const { user, routeId, setRoute } = useAuth();

  const [dialog, setDialog] = useState<DialogKind>(null);
  const [draftId, setDraftId] = useState<string | null>(routeId);
  const [routeError, setRouteError] = useState('');
  const [feedback, setFeedback] = useState('');

  const route = routeById(routeId);

  if (!user) {
    return null; // ProtectedRoute не пропустит сюда без входа
  }

  const openRouteDialog = () => {
    setDraftId(routeId);
    setRouteError('');
    setDialog('route');
  };

  const confirmRoute = () => {
    const chosen = routeById(draftId);
    if (!chosen) {
      setRouteError(CABINET_COPY.routeDialog.requiredError);
      return;
    }
    setRoute(chosen.id);
    setRouteError('');
    setFeedback(CABINET_COPY.routeDialog.saved);
    setDialog(null);
  };

  const confirmRemove = () => {
    setRoute(null);
    setFeedback(CABINET_COPY.routeDialog.removed);
    setDialog(null);
  };

  const openSample = () => {
    if (route) {
      setDialog('sample');
    } else {
      setFeedback(CABINET_COPY.preparation.emptyNote);
    }
  };

  return (
    <>
      <div className="page-intro">
        <div>
          <h1>{CABINET_COPY.page.journeyTitle}</h1>
          <p className="page-context">{CABINET_COPY.page.context}</p>
        </div>
        <p className="participant">
          {CABINET_COPY.page.participantLabel} <br />
          <strong>{user.login}</strong>
        </p>
      </div>

      <div id="cabinet-content" role="tabpanel" aria-labelledby="tab-journey">
        {route ? (
          <section
            className="route-stage"
            aria-labelledby="route-title"
            data-route-stage=""
            data-route-id={route.id}
          >
            <div className="route-stage__metric">
              <FlowerMark />
              <div className="distance-text">
                <span>{route.distance.replace(' км', '')}</span>
                <small>км</small>
              </div>
            </div>
            <div className="route-stage__body">
              <h2 className="route-title" id="route-title">
                {route.title}
              </h2>
              <div className="route-kicker">{route.kind}</div>
              <p className="route-description">{route.description}</p>
              <p className="status-label">{CABINET_COPY.page.applicationStatus}</p>
              <div className="route-actions">
                <button className="primary-button" type="button" onClick={openRouteDialog}>
                  {CABINET_COPY.actions.changeRoute}
                </button>
                <button
                  className="plain-button"
                  type="button"
                  onClick={() => {
                    setFeedback('');
                    setDialog('remove');
                  }}
                >
                  {CABINET_COPY.actions.removeRoute}
                </button>
              </div>
            </div>
          </section>
        ) : (
          <section
            className="route-stage empty-stage"
            aria-labelledby="route-title"
            data-route-stage=""
          >
            <div className="route-stage__metric">
              <FlowerMark />
              <span className="empty-distance" aria-hidden="true">
                ?
              </span>
            </div>
            <div className="route-stage__body">
              <h2 className="route-title" id="route-title">
                {CABINET_COPY.page.emptyTitle}
              </h2>
              <p className="route-description">{CABINET_COPY.page.emptyBody}</p>
              <div className="route-actions">
                <button className="primary-button" type="button" onClick={openRouteDialog}>
                  {CABINET_COPY.actions.chooseRoute}
                </button>
              </div>
            </div>
          </section>
        )}

        <section className="preparation" aria-labelledby="preparation-title">
          <div>
            <h2 id="preparation-title">{CABINET_COPY.preparation.title}</h2>
          </div>
          <div className="preparation-body">
            <p>{CABINET_COPY.preparation.body}</p>
            <div className="prep-actions">
              <button
                className="secondary-button"
                type="button"
                onClick={openSample}
                disabled={!route}
              >
                {CABINET_COPY.actions.sample}
              </button>
              <a
                href={CABINET_COPY.externalLinks.organizerCabinet}
                target="_blank"
                rel="noopener noreferrer"
              >
                {CABINET_COPY.actions.organizerCabinet}
              </a>
            </div>
            {!route ? (
              <p className="disabled-note">{CABINET_COPY.preparation.emptyNote}</p>
            ) : null}
            <div className="service-links">
              <a
                href={CABINET_COPY.externalLinks.participants}
                target="_blank"
                rel="noopener noreferrer"
              >
                Как участвовать
              </a>
              <a
                href={CABINET_COPY.externalLinks.routes}
                target="_blank"
                rel="noopener noreferrer"
              >
                Все маршруты
              </a>
              <a
                href={CABINET_COPY.externalLinks.faq}
                target="_blank"
                rel="noopener noreferrer"
              >
                Частые вопросы
              </a>
            </div>
          </div>
        </section>

        <p className="form-feedback" role="status">
          {feedback}
        </p>
      </div>

      <RouteDialog
        open={dialog === 'route'}
        draftId={draftId}
        error={routeError}
        onSelect={(nextId) => {
          setDraftId(nextId);
          setRouteError('');
        }}
        onConfirm={confirmRoute}
        onCancel={() => {
          setDraftId(routeId);
          setRouteError('');
          setDialog(null);
        }}
      />
      <RemoveRouteDialog
        open={dialog === 'remove'}
        onConfirm={confirmRemove}
        onCancel={() => setDialog(null)}
      />
      <SampleDialog
        open={dialog === 'sample'}
        login={user.login}
        email={user.email}
        routeId={routeId}
        onClose={() => setDialog(null)}
      />
    </>
  );
}
