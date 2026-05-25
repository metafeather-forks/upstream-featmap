import { Roles } from "../../core/misc";


export interface IApplication {
    mode: string
    workspaces: IWorkspace[]
    memberships: IMembership[]
    account?: IAccount
    messages: IMessage[]
}

export interface IMembership {
    id: string
    workspaceId: string
    accountId: string
    level: Roles
    name: string
    email: string
    createdAt: string
}

export interface IWorkspace {
    id: string
    name: string
    createdAt: string
    allowExternalSharing: boolean
    euVat: string
    status: string
}

export interface IAccount {
    id: string
    name: string
    email: string
    createdAt: string
    emailConfirmed: boolean
    emailConfirmationSentTo: string
    emailConfirmationPending: boolean
}

export interface IInvite {
    id: string
    workspaceId: string
    email: string
    level: string
    code: string
    createdBy: string
    createdByName: string
    createdAt: string
    createdByEmail: string
    workspaceName: string
}

export type messageTypes = "success" | "fail"

export interface IMessage { id: string, type: messageTypes, message: string }
