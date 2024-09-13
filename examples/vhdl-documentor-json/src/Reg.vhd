library IEEE;

use IEEE.std_logic_1164.all;

entity Reg is

  port(
    iCLK : in  std_logic;
    iD   : in  integer;
    oQ   : out integer
    );

end Reg;

architecture behavior of Reg is
begin

  process(iCLK, iD)
  begin
    if rising_edge(iCLK) then
      oQ <= iD;
    end if;
  end process;

end behavior;
