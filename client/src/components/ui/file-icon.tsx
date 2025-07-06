import { MdAttachFile } from "react-icons/md";
import { IoDocumentTextOutline } from "react-icons/io5";
import { FaRegImage } from "react-icons/fa6";

export default function FileIcon({ ext }: { ext: string }) {
  if (ext.toLocaleLowerCase() === "pdf") {
    return <IoDocumentTextOutline />;
  }

  if (
    ext.toLocaleLowerCase() === "jpg" ||
    ext.toLocaleLowerCase() === "jpeg" ||
    ext.toLocaleLowerCase() === "png" ||
    ext.toLocaleLowerCase() === "gif" ||
    ext.toLocaleLowerCase() === "webp"
  ) {
    return <FaRegImage />;
  }

  return <MdAttachFile />;
}
