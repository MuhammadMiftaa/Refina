import Template from "../../components/layouts/template";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";
import { RadioGroup, RadioGroupItem } from "../../components/ui/radio-group";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../../components/ui/select";

export default function FormAddPage() {
  return (
    <Template>
      <div className="w-full min-h-screen bg-secondary py-5 px-3">
        <div className="font-poppins text-center">
          <h1 className="text-2xl font-bold text-white">
            Add your transaction here.
          </h1>
          <p className="text-sm font-light text-zinc-400">
            Keep track of your income and expenses.
          </p>
        </div>

        <form className="font-poppins py-6 px-2 flex flex-col gap-5">
          <div className="grid w-full max-w-sm items-center gap-1.5">
            <Label
              className="text-2xl font-semibold text-white"
              htmlFor="description"
            >
              Description
            </Label>
            <Input
              className="shadow-zinc-400 border-zinc-400 text-zinc-200 text-xl"
              type="text"
              id="description"
              placeholder="Belanja Bulanan"
            />
          </div>
          <div className="grid w-full max-w-sm items-center gap-1.5">
            <Label
              className="text-2xl font-semibold text-white"
              htmlFor="category"
            >
              Category
            </Label>
            <Select>
              <SelectTrigger className="shadow-zinc-400 border-zinc-400 text-zinc-500 focus:ring-0">
                <SelectValue placeholder="Select a timezone" />
              </SelectTrigger>
              <SelectContent className="font-poppins bg-secondary text-zinc-300">
                <SelectGroup>
                  <SelectLabel className="text-white">
                    North America
                  </SelectLabel>
                  <SelectItem value="est">
                    Eastern Standard Time (EST)
                  </SelectItem>
                  <SelectItem value="cst">
                    Central Standard Time (CST)
                  </SelectItem>
                  <SelectItem value="mst">
                    Mountain Standard Time (MST)
                  </SelectItem>
                  <SelectItem value="pst">
                    Pacific Standard Time (PST)
                  </SelectItem>
                  <SelectItem value="akst">
                    Alaska Standard Time (AKST)
                  </SelectItem>
                  <SelectItem value="hst">
                    Hawaii Standard Time (HST)
                  </SelectItem>
                </SelectGroup>
                <SelectGroup>
                  <SelectLabel className="text-white">
                    Europe & Africa
                  </SelectLabel>
                  <SelectItem value="gmt">Greenwich Mean Time (GMT)</SelectItem>
                  <SelectItem value="cet">
                    Central European Time (CET)
                  </SelectItem>
                  <SelectItem value="eet">
                    Eastern European Time (EET)
                  </SelectItem>
                  <SelectItem value="west">
                    Western European Summer Time (WEST)
                  </SelectItem>
                  <SelectItem value="cat">Central Africa Time (CAT)</SelectItem>
                  <SelectItem value="eat">East Africa Time (EAT)</SelectItem>
                </SelectGroup>
                <SelectGroup>
                  <SelectLabel className="text-white">Asia</SelectLabel>
                  <SelectItem value="msk">Moscow Time (MSK)</SelectItem>
                  <SelectItem value="ist">India Standard Time (IST)</SelectItem>
                  <SelectItem value="cst_china">
                    China Standard Time (CST)
                  </SelectItem>
                  <SelectItem value="jst">Japan Standard Time (JST)</SelectItem>
                  <SelectItem value="kst">Korea Standard Time (KST)</SelectItem>
                  <SelectItem value="ist_indonesia">
                    Indonesia Central Standard Time (WITA)
                  </SelectItem>
                </SelectGroup>
                <SelectGroup>
                  <SelectLabel className="text-white">
                    Australia & Pacific
                  </SelectLabel>
                  <SelectItem value="awst">
                    Australian Western Standard Time (AWST)
                  </SelectItem>
                  <SelectItem value="acst">
                    Australian Central Standard Time (ACST)
                  </SelectItem>
                  <SelectItem value="aest">
                    Australian Eastern Standard Time (AEST)
                  </SelectItem>
                  <SelectItem value="nzst">
                    New Zealand Standard Time (NZST)
                  </SelectItem>
                  <SelectItem value="fjt">Fiji Time (FJT)</SelectItem>
                </SelectGroup>
                <SelectGroup>
                  <SelectLabel className="text-white">
                    South America
                  </SelectLabel>
                  <SelectItem value="art">Argentina Time (ART)</SelectItem>
                  <SelectItem value="bot">Bolivia Time (BOT)</SelectItem>
                  <SelectItem value="brt">Brasilia Time (BRT)</SelectItem>
                  <SelectItem value="clt">Chile Standard Time (CLT)</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          <div className="grid w-full max-w-sm items-center gap-1.5 relative">
            <Label
              className="text-2xl font-semibold text-white"
              htmlFor="amount"
            >
              Amount
            </Label>
            <Input
              className="shadow-zinc-400 border-zinc-400 text-zinc-200 text-xl"
              type="number"
              id="amount"
              placeholder="500.000"
            />
            <div className="h-7 w-4 absolute right-3 bottom-1 bg-secondary"></div>
          </div>
          <div className="grid w-full max-w-sm items-center gap-1.5 relative">
            <Label className="text-2xl font-semibold text-white" htmlFor="type">
              Transacion Type
            </Label>
            <RadioGroup className="flex gap-6" defaultValue="income">
              <div className="flex items-center py-1 px-4 rounded text-white border border-green-500 shadow-sm shadow-green-500 cursor-pointer">
                <RadioGroupItem
                  value="income"
                  className="scale-100 hidden"
                  id="income"
                />
                <Label className="text-lg cursor-pointer" htmlFor="income">Income</Label>
              </div>
              <div className="flex items-center py-1 px-4 rounded text-white border border-red-500 shadow-sm shadow-red-500 cursor-pointer">
                <RadioGroupItem
                  value="expense"
                  className="scale-100 hidden"
                  id="expense"
                />
                <Label className="text-lg cursor-pointer" htmlFor="expense">Expense</Label>
              </div>
            </RadioGroup>
          </div>
        </form>
      </div>
    </Template>
  );
}
