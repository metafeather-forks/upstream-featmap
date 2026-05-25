import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Formik, FormikHelpers as FormikActions, FormikProps, Form } from 'formik';
import { API_ACCEPT_INVITE, API_GET_INVITE } from '../api'
import { Link } from 'react-router-dom'
import { Button } from '../components/elements';
import { memberLevelToTitle } from "../core/misc";
import { IInvite } from '../store/application/types';

function AcceptInvitePage() {
  const { code } = useParams<{ code: string }>();
  const [invite, setInvite] = useState<IInvite | undefined>();
  const [notFound, setNotFound] = useState(false);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    API_GET_INVITE(code).then(response => {
      if (response.ok) {
        response.json().then((data: IInvite) => setInvite(data))
      } else {
        setNotFound(true)
      }
    })
  }, [code]);

  if (notFound) {
    return <div className="p-2">Invitation not found. Back to <Link className="link" to="/">Featmap</Link>.</div>
  }
  if (success) {
    return <div className="p-2">You are now a member of <b>{invite!.workspaceName}</b>! Back to <Link className="link" to="/">Featmap</Link>.</div>
  }
  if (invite) {
    return (
      <div className="flex p-2 w-full justify-center items-center flex-col">
        <div className="flex p-3 max-w-xl w-full items-center flex-col">
          <div className="flex p-2 flex-col items-baseline">
            <div className="p-1"><h1 className="text-3xl font-medium">Invitation to join workspace</h1></div>
          </div>
          <p>The workspace name is <b>{invite.workspaceName}</b> and you will join as <b>{memberLevelToTitle(invite.level)}</b>.</p>
          <div>
            <Formik initialValues={{}} onSubmit={(values: {}, actions: FormikActions<{}>) => {
              actions.setStatus("")
              API_ACCEPT_INVITE(code).then(response => {
                if (response.ok) {
                  setSuccess(true)
                } else {
                  response.json().then((data: { message: string }) => actions.setStatus(data.message))
                }
              })
              actions.setSubmitting(false)
            }}>
              {(formikBag: FormikProps<{}>) => (
                <Form>
                  <div className="p-2 font-bold text-red">{formikBag.status}</div>
                  <div className="flex w-full text-lg justify-center">
                    <div><Button title="Accept invite" primary submit /></div>
                  </div>
                </Form>
              )}
            </Formik>
          </div>
          <div className="flex p-2 flex-col">
            <div className="p-1 text-center">Not a member? <Link target="_blank" className="link" to="/account/signup">Create an account</Link></div>
          </div>
        </div>
      </div>
    )
  }
  return <div>Loading</div>;
}

export default AcceptInvitePage;
