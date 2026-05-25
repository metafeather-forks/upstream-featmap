import { Link } from 'react-router-dom'
import { Color, colorToBorderColorClass, dbAnnotationsFromNames } from '../core/misc';

type Props = {
  title: string
  link: string
  status?: string
  small?: boolean
  color?: Color
  bottomLink?: () => void
  nbrOfItems?: number
  nbrOfComments?: number
  annotations: string
  estimate?: number
};

function Card(props: Props) {
  const annos = dbAnnotationsFromNames(props.annotations)
  const color = props.color && (props.color !== Color.WHITE) ? props.color : null

  return (
    <div>
      <div style={{ fontSize: '12px' }} className={"flex flex-row flex-no-shrink w-36 rounded-sm bg-white overflow-hidden border " + (props.small ? " " : " h-24 ") + (color ? " border-l-4 " + colorToBorderColorClass(color) + " " : " " + colorToBorderColorClass(Color.WHITE) + " ")}>
        <Link className="flex flex-col flex-grow" to={props.link}>
          <div className="flex-grow p-2 font-normal overflow-hidden">
            <span className={props.status === "CLOSED" ? "line-through" : ""}> {props.title} </span>
          </div>
          <div className="flex p-2 flex-row">
            {props.nbrOfItems && !(props.nbrOfItems === 0) ?
              <div className="flex"><div>{props.nbrOfItems} items</div></div> : null}
            {props.nbrOfComments! > 0 ?
              <div title={props.nbrOfComments + " comments"} className="whitespace-nowrap">
                {props.nbrOfComments!}<i style={{ fontSize: "12px" }} className="material-icons align-middle">chat_bubble_outline</i>
              </div> : null}
            <div className="flex flex-grow" />
            {props.annotations ?
              annos.annotations.map((a, i) =>
                <div className="bg-gray-200" key={i}><i style={{ fontSize: "14px" }} title={a.description} className="text-gray-700 material-icons align-middle">{a.icon}</i></div>
              ) : null}
            {props.estimate && !(props.estimate === 0) ?
              <div className="whitespace-nowrap bg-gray-200 px-1"><span title={"Size is " + props.estimate}> {props.estimate}</span></div> : null}
          </div>
        </Link>
      </div>
      {props.bottomLink &&
        <div className="flex w-36 h-6 flex-no-shrink -mt-1 -mb-2 justify-center">
          <div className="flex showme font-bold text-xl"><button className="hover:text-gray-800 text-gray-500" onClick={props.bottomLink}>+</button></div>
        </div>}
    </div>
  );
}

export default Card;
