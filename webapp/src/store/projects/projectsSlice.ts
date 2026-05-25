import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { IProject } from './types'
import { AppState } from '..'
import { createSelector } from 'reselect'

export { IProject } from './types'

interface ProjectsState {
  items: IProject[]
}

const initialState: ProjectsState = {
  items: [],
}

const projectsSlice = createSlice({
  name: 'projects',
  initialState,
  reducers: {
    createProject(state, action: PayloadAction<IProject>) {
      action.payload.kind = 'project'
      state.items.push(action.payload)
    },
    loadProjects(state, action: PayloadAction<IProject[]>) {
      state.items = action.payload.map(x => ({ ...x, kind: 'project' } as IProject))
    },
    updateProject(state, action: PayloadAction<IProject>) {
      action.payload.kind = 'project'
      const idx = state.items.findIndex(x => x.id === action.payload.id)
      if (idx >= 0) state.items[idx] = action.payload
    },
    deleteProject(state, action: PayloadAction<string>) {
      state.items = state.items.filter(x => x.id !== action.payload)
    },
  },
})

export const { createProject, loadProjects, updateProject, deleteProject } = projectsSlice.actions

// Selectors
const getProjectsState = (state: AppState) => state.projects

export const projectsSelector = createSelector([getProjectsState], s =>
  [...s.items].sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
)

export const getProjectById = (pp: IProject[], id: string) => pp.find(x => x.id === id)
export const sortProjectsByCreateDate = (pp: IProject[]): IProject[] =>
  [...pp].sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

export default projectsSlice.reducer
