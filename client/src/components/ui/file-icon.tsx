import { MdAttachFile } from "react-icons/md";
import { IoDocumentTextOutline } from "react-icons/io5";
import { FaRegImage } from "react-icons/fa6";

export default function FileIcon({ ext }: { ext: string }) {
  if (ext.toLowerCase() === "pdf") {
    return <IoDocumentTextOutline />;
  }

  if (
    ext.toLowerCase() === "jpg" ||
    ext.toLowerCase() === "jpeg" ||
    ext.toLowerCase() === "png" ||
    ext.toLowerCase() === "gif" ||
    ext.toLowerCase() === "webp"
  ) {
    return <FaRegImage />;
  }

  return <MdAttachFile />;
}
