# Code Context

## Files Retrieved
1. `webapp/package.json` - Current dependencies: React 18, Redux 4, react-redux 7, React Router 5, react-beautiful-dnd 13, formik 2, yup 0.32, uuid 8, react-onclickoutside 6, react-markdown 6, typesafe-actions 5, reselect 4
2. `webapp/tsconfig.json` - target es2020, strict:true, jsx:react-jsx, moduleResolution:bundler
3. `webapp/vite.config.ts` - Vite 5, builds to `build/`, dev on port 3000
4. `webapp/index.html` - Loads Tailwind from CDN, Material Icons CDN, Google Fonts CDN, Stripe.js CDN
5. `webapp/src/index.tsx` (lines 1-15) - Uses deprecated `ReactDOM.render` instead of `createRoot`
6. `webapp/src/App.tsx` (lines 1-52) - Class component with React Router v5 `<Switch>` and `connect()`
7. `webapp/src/configureStore.ts` (lines 1-11) - Manual `createStore` with deprecated Redux DevTools extension pattern
8. `webapp/src/store/index.ts` (lines 1-47) - 9 domain reducers combined with `combineReducers`, hand-written `AllActions` union type
9. `webapp/src/api/index.ts` (~800 lines) - ~80 raw `fetch()` API functions with duplicated `response.json().then()` patterns
10. `webapp/src/components/elements.tsx` - Uses deprecated `SFC` type, no `children` in `CardLayout` props
11. `webapp/src/components/Board.tsx` (842 lines) - Largest component, class-based, react-beautiful-dnd, inline API calls, 8 useState equivalents as `this.state`
12. `webapp/src/components/EntityDetailsBody.tsx` (725 lines) - Class component, ~30 `alert()` error calls, inline API calls, `window.confirm()`
13. `webapp/src/components/EntityDetailsTitle.tsx` (319 lines) - Formik render-prop, react-onclickoutside HOC
14. `webapp/src/store/application/actions.ts` - Manual action types enum, typesafe-actions, dispatch thunks
15. `webapp/src/store/features/reducers.ts` - Switch-case reducer with manual state spreading
16. `webapp/src/store/features/selectors.ts` - reselect createSelector usage
17. `webapp/src/core/misc.ts` - Color/misc enums, annotations, ~250 lines of country data
18. `webapp/src/core/card.ts` - EntityTypes union type for polymorphic card handling

## Key Code

### Current Redux Pattern (all 9 stores follow this)
```typescript
// actions.ts
export enum ActionTypes { CREATE_FEATURE = 'CREATE_FEATURE', ... }
export interface createFeature { type: ActionTypes.CREATE_FEATURE, payload: IFeature }
export const createFeatureAction = (x: IFeature) => action(ActionTypes.CREATE_FEATURE, x)
// reducers.ts
export function reducer(state = initialState, action: Actions) {
  switch (action.type) {
    case ActionTypes.CREATE_FEATURE: return { ...state, items: [...state.items, action.payload] }
  }
}
```

### Current Component Pattern (38 class components)
```typescript
// Every component uses this pattern:
class X extends Component<Props, State> {
  constructor(props) { super(props); this.state = { ... } }
  componentDidMount() { API_FETCH_X().then(r => r.json().then(data => ...)) }
  render() { return (<div>...</div>) }
}
export default connect(mapStateToProps, mapDispatchToProps)(X);
```

### Current API Pattern
```typescript
// ~80 functions in api/index.ts, all raw fetch()
export const API_RENAME_FEATURE = async (workspaceId, id, title) =>
  await fetch(endpoint + "/features/" + id + "/rename", {
    method: 'POST', headers: { 'Workspace': workspaceId }, credentials: 'include',
    body: JSON.stringify({ title })
  });
// Components call these inline with .then().catch(), ~40 alert() error handlers
```

### Current Routing
```typescript
// React Router v5
<Switch>
  <Route exact path="/account/login" component={LoginPage} />
  <Route path="/" component={IndexPage} />
</Switch>
// Pages access: this.props.match.params, this.props.history.push()
```

## Architecture

### Directory Structure
```
webapp/src/
  api/index.ts          -- ~80 raw fetch functions, no service layer
  App.tsx, App.css      -- root layout, React Router v5 Switch
  index.tsx             -- entry: ReactDOM.render + Provider + BrowserRouter
  configureStore.ts     -- manual createStore with compose
  store/
    index.ts            -- combineReducers, AllActions type
    application/        -- app state, messages system
    features/           -- feature cards CRUD + drag move
    milestones/         -- milestone columns CRUD + reorder
    projects/           -- project CRUD
    workflows/          -- workflow rows CRUD + reorder
    subworkflows/       -- subworkflow groups within workflows
    personas/           -- persona management
    workflowpersonas/   -- persona-workflow assignments
    featurecomments/    -- comments on features
  components/
    Board.tsx           -- main board: DnD context, milestone/workflow/subworkflow/feature grid
    Card.tsx            -- feature card display (Link wrapper)
    EntityDetailsBody.tsx -- entity detail: CRUD, color, annotations, estimate, delete
    EntityDetailsTitle.tsx -- inline editable title (Formik)
    EntityDetailsDescription.tsx -- inline editable description (Formik)
    EntityDetailsComments.tsx -- comment list + new comment form (Formik)
    EntityDetailsAnnotations.tsx -- annotation tags (Formik)
    EntityDetailsModal.tsx -- modal wrapper
    CreateCardModal.tsx -- create feature/milestone/workflow/subworkflow (Formik)
    CreateProjectModal.tsx -- create project (Formik)
    CreateWorkspaceModal.tsx -- create workspace (Formik)
    Personas.tsx       -- persona bar: create/edit/assign (Formik)
    ContextMenu.tsx     -- generic dropdown menu (onClickOutside HOC)
    Header.tsx, Footer.tsx, Messages.tsx -- layout chrome
    elements.tsx        -- Button, CardLayout primitives
    Comment.tsx         -- single comment display + inline edit
    NewCard.tsx, NewDimCard.tsx, EmptyCard.tsx -- simple display components
  pages/
    IndexPage.tsx       -- app loader, auth gate, sub-routing
    ProjectPage.tsx     -- project detail, loads all entities, renders Board
    WorkspacePage.tsx   -- workspace with project list + sub-routing
    ...auth pages...    -- Login, SignUp, ResetPassword, VerifyEmail, Logout, AcceptInvite
    ...settings pages...-- AccountPage, WorkspaceSettingsPage
  core/
    lexorank.ts, lexorank.test.tsx -- lexicographic ranking for drag reorder
    misc.ts            -- colors, annotations, countries, role helpers
    card.ts            -- EntityTypes union type
```

### Data Flow
1. `IndexPage` fetches `/account/app` -> dispatches `receiveAppAction` -> loads workspaces/memberships/account
2. `WorkspacePage` fetches projects for workspace -> dispatches `loadProjectsAction`
3. `ProjectPage` fetches all project entities (milestones, workflows, features, comments, personas) -> dispatches load actions for each domain
4. `Board` renders the grid, handles drag-and-drop with optimistic Redux dispatches + async API calls
5. Entity detail modals do inline API calls with `alert()` on failure, dispatch `updateXAction` on success
6. Forms use Formik render-prop pattern with Yup validation schemas

## Start Here

Open `webapp/TODO.md` for the full modernization plan. The recommended first action is Phase 1.1 (bump Vite to 6) combined with Phase 1.2 (createRoot migration in `src/index.tsx`), as these are low-risk and unlock everything else.

The most impactful single task is Phase 2.3 (convert Redux stores to RTK createSlice), which eliminates ~500 lines of boilerplate across 9 domains and makes Phase 3 (component migration) cleaner by enabling `useAppDispatch`/`useAppSelector` hooks.
